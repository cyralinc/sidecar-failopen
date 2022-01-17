package logging

import (
	"os"

	_ "github.com/cyralinc/sidecar-failopen/internal/config"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

func Info(template string, args ...interface{}) {
	defer logger.Sync()
	logger.Infof(template, args...)
}

func Error(template string, args ...interface{}) {
	defer logger.Sync()
	logger.Errorf(template, args...)
}

func Panic(template string, args ...interface{}) {
	defer logger.Sync()
	logger.Panicf(template, args...)
}

func Warn(template string, args ...interface{}) {
	defer logger.Sync()
	logger.Warnf(template, args...)
}

func Debug(template string, args ...interface{}) {
	defer logger.Sync()
	logger.Debugf(template, args...)
}

func init() {
	logLevel := viper.GetString("log_level")
	atom := zap.NewAtomicLevel()

	encoderCfg := zap.NewProductionEncoderConfig()

	logger = zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	)).Sugar()

	atom.UnmarshalText([]byte(logLevel))

}
