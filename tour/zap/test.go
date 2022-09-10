package zap

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

var lg *zap.Logger

func init() {
	str, _ := os.Getwd()
	logPath := filepath.Join(str, fmt.Sprintf("test-log-%v.log", time.Now().Format("2006-01-02")))
	zapLogPath := filepath.Join("winfile:///", logPath)
	defer os.RemoveAll(logPath)
	lcfg := zap.NewProductionConfig()
	zap.RegisterSink("winfile", newWinFileSink)
	lcfg.OutputPaths = []string{zapLogPath}
	lcfg.ErrorOutputPaths = []string{zapLogPath}
	//lcfg.Encoding ="console"
	//lcfg.EncoderConfig = zapcore.EncoderConfig{
	//	MessageKey:       "message",
	//	LevelKey:         "level",
	//	TimeKey:          "",
	//	NameKey:          "",
	//	CallerKey:        "",
	//	FunctionKey:      "",
	//	StacktraceKey:    "",
	//	LineEnding:       "",
	//	EncodeLevel:      nil,
	//	EncodeTime:       nil,
	//	EncodeDuration:   nil,
	//	EncodeCaller:     nil,
	//	EncodeName:       nil,
	//	ConsoleSeparator: "",
	//}

	var err error
	lg, err = lcfg.Build()
	if err != nil {
		log.Fatal(err)
	}
	lg.Sync()
}
func WriteLog(s string) {
	lg.Info(s)

	//注意！ fatal 会直接导致程序退出  那么就不会接着打印日志了！！！
	//lg.Fatal("helloworld")

	// For some users, the presets offered by the NewProduction, NewDevelopment,
	// and NewExample constructors won't be appropriate. For most of those
	// users, the bundled Config struct offers the right balance of flexibility
	// and convenience. (For more complex needs, see the AdvancedConfiguration
	// example.)
	//
	// See the documentation for Config and zapcore.EncoderConfig for all the
	// available options.
	//rawJSON := []byte(`{
	//  "level": "debug",
	//  "encoding": "json",
	//  "outputPaths": ["stdout", "/tmp/logs"],
	//  "errorOutputPaths": ["stderr"],
	//  "initialFields": {"foo": "bar"},
	//  "encoderConfig": {
	//    "messageKey": "message",
	//    "levelKey": "level",
	//    "levelEncoder": "lowercase"
	//  }
	//}`)

	//var cfg zap.Config
	//if err := json.Unmarshal(rawJSON, &cfg); err != nil {
	//	panic(err)
	//}
	//logger, err = cfg.Build()
	//if err != nil {
	//	panic(err)
	//}
	//defer logger.Sync()
	//
	//logger.Info("logger construction succeeded")
}

func newWinFileSink(u *url.URL) (zap.Sink, error) {
	// Remove leading slash left by url.Parse()
	var name string
	if u.Path != "" {
		name = u.Path[1:]
	} else if u.Opaque != "" {
		name = u.Opaque[1:]
	} else {
		return nil, errors.New("path error")
	}
	return os.OpenFile(name, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
}
