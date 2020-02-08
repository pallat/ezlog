package ezlog_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/pallat/ezlog"
	"github.com/sirupsen/logrus"
)

type LogStruct struct {
	Latency string `json:"latency"`
	Level   string `json:"level"`
	Msg     string `json:"msg"`
	Name    string `json:"name"`
	Time    string `json:"time"`
}

func TestLogJSONFormat(t *testing.T) {
	var want map[string]string
	json.Unmarshal([]byte(`{"latency":"1s","level":"info","msg":"","name":"api","time":"2020-02-08T12:44:20+07:00"}`), &want)

	given := fmt.Sprintf("latency=%s name=api", time.Second)

	getBuf := new(bytes.Buffer)
	ezlog.DefaultLogger.SetOutput(getBuf)
	ezlog.Print(given)

	var get map[string]string
	json.Unmarshal(getBuf.Bytes(), &get)

	if _, ok := get["latency"]; !ok {
		t.Error("not found the key latency we want")
		return
	}

	if _, ok := get["name"]; !ok {
		t.Error("not found the key name we want")
		return
	}

	if want["latency"] != get["latency"] {
		t.Errorf("%s is wanted but get %s\n", want["latency"], get["latency"])
	}
	if want["name"] != get["name"] {
		t.Errorf("%s is wanted but get %s\n", want["name"], get["name"])
	}
}

func TestLogTextFormat(t *testing.T) {
	ezlog.DefaultLogger.Formatter = &logrus.TextFormatter{}

	given := fmt.Sprintf("latency=%s name=api", time.Second)

	getBuf := new(bytes.Buffer)
	ezlog.DefaultLogger.SetOutput(getBuf)
	ezlog.Print(given)

	if !strings.Contains(getBuf.String(), "latency=1s name=api") {
		t.Errorf("it should contains %s in log", getBuf)
	}
}
