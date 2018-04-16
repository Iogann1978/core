package logger

import (
	"bytes"
	"encoding/json"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	"testing"
)

type (
	logMessage struct {
		Level     string                 `json:"level"`
		Timestamp float64                `json:"ts"`
		Message   string                 `json:"msg"`
		Args      map[string]interface{} `json:"args"`
	}
)

var (
	testMapData = map[string]interface{}{`str`: `123`, `number`: 123.0}
)

func getLogger() (Logger, *bytes.Buffer) {
	b := new(bytes.Buffer)
	bs := zapcore.AddSync(b)
	l := NewZapLogger(bs)
	return l, b
}

func TestNewZapLogger(t *testing.T) {
	var m logMessage
	l, b := getLogger()
	l.Warn(`init`, map[string]interface{}{``: ``})
	if err := json.Unmarshal(b.Bytes(), &m); err != nil {
		t.Fatal(err)
	} else {
		assert.Equal(t, `init`, m.Message)
		assert.Equal(t, `warn`, m.Level)
	}
}

func TestZapLogger_Debug(t *testing.T) {
	var m logMessage
	l, b := getLogger()
	l.Debug(`debug`, testMapData)
	if err := json.Unmarshal(b.Bytes(), &m); err != nil {
		t.Fatal(err)
	} else {
		assert.Equal(t, `debug`, m.Level)
		assert.Equal(t, `debug`, m.Message)
		if diff := deep.Equal(testMapData, m.Args); diff != nil {
			t.Error(diff)
		}
	}
}

func TestZapLogger_Info(t *testing.T) {
	var m logMessage
	l, b := getLogger()
	l.Info(`info`, testMapData)
	if err := json.Unmarshal(b.Bytes(), &m); err != nil {
		t.Fatal(err)
	} else {
		assert.Equal(t, `info`, m.Level)
		assert.Equal(t, `info`, m.Message)
		if diff := deep.Equal(testMapData, m.Args); diff != nil {
			t.Error(diff)
		}
	}
}

func TestZapLogger_Panic(t *testing.T) {
	var m logMessage
	l, b := getLogger()

	defer func(t *testing.T) {
		if r := recover(); r != nil {
			assert.Equal(t, `panic`, r.(string))
			if err := json.Unmarshal(b.Bytes(), &m); err != nil {
				t.Fatal(err)
			} else {
				assert.Equal(t, `panic`, m.Level)
				assert.Equal(t, `panic`, m.Message)
				if diff := deep.Equal(testMapData, m.Args); diff != nil {
					t.Error(diff)
				}
			}
		} else {
			t.Fatal(`recover is nil`)
		}
	}(t)
	l.Panic(`panic`, testMapData)

}

func TestZapLogger_Warn(t *testing.T) {
	var m logMessage
	l, b := getLogger()
	l.Warn(`warn`, testMapData)
	if err := json.Unmarshal(b.Bytes(), &m); err != nil {
		t.Fatal(err)
	} else {
		assert.Equal(t, `warn`, m.Level)
		assert.Equal(t, `warn`, m.Message)
		if diff := deep.Equal(testMapData, m.Args); diff != nil {
			t.Error(diff)
		}
	}
}
