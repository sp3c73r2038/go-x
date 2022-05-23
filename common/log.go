package common

import (
	"encoding/json"
	"fmt"

	"github.com/kr/pretty"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var LoggerConfig = zap.NewDevelopmentConfig()
var RawLogger, _ = LoggerConfig.Build()
var Logger = RawLogger.Sugar()

func Must(err error) {
	if err != nil {
		Logger.Fatal(err)
	}
}

func SetupLogger() {
	// LoggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	RawLogger, _ = LoggerConfig.Build()
	Logger = RawLogger.Sugar()
}

func Pretty(obj interface{}) string {
	return fmt.Sprintf("%# v", pretty.Formatter(obj))
}

func Dump(obj interface{}) string {
	var err error
	var b []byte
	b, err = json.Marshal(obj)
	Must(err)
	return string(b)
}

var LevelInfo = zapcore.InfoLevel
var LevelDebug = zapcore.DebugLevel
