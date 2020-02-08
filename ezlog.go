package ezlog

import (
	"bytes"
	"fmt"

	"github.com/sirupsen/logrus"
)

var DefaultLogger = logrus.New()

func init() {
	DefaultLogger.Formatter = &logrus.JSONFormatter{}
}

func Print(args ...interface{}) {
	s := fmt.Sprint(args...)
	all := logrus.Fields{}
	m := ParseLogfmt(s)
	for k, v := range m {
		all[k] = v
	}

	DefaultLogger.WithFields(all).Info()

}

func Printf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	all := logrus.Fields{}
	m := ParseLogfmt(s)
	for k, v := range m {
		all[k] = v
	}

	DefaultLogger.WithFields(all).Info()
}

// ParseLogfmt from https://gist.github.com/alexisvisco/4b846978c9346e4eaf618bb632c0693a
func ParseLogfmt(msg string) map[string]string {

	type kv struct {
		key, val string
	}

	var pair *kv = nil
	pairs := make(map[string]string)
	buf := bytes.NewBuffer([]byte{})

	var (
		escape  = false
		garbage = false
		quoted  = false
	)

	completePair := func(buffer *bytes.Buffer, pair *kv) kv {
		if pair != nil {
			return kv{pair.key, buffer.String()}
		} else {
			return kv{buffer.String(), ""}
		}
	}

	for _, c := range msg {
		if !quoted && c == ' ' {
			if buf.Len() != 0 {
				if !garbage {
					p := completePair(buf, pair)
					pairs[p.key] = p.val
					pair = nil
				}
				buf.Reset()
			}
			garbage = false
		} else if !quoted && c == '=' {
			if buf.Len() != 0 {
				pair = &kv{key: buf.String(), val: ""}
				buf.Reset()
			} else {
				garbage = true
			}
		} else if quoted && c == '\\' {
			escape = true
		} else if c == '"' {
			if escape {
				buf.WriteRune(c)
				escape = false
			} else {
				quoted = !quoted
			}
		} else {
			if escape {
				buf.WriteRune('\\')
				escape = false
			}
			buf.WriteRune(c)
		}
	}

	if !garbage {
		p := completePair(buf, pair)
		pairs[p.key] = p.val
	}

	return pairs
}
