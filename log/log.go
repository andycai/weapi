package log

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger
var sugar *zap.SugaredLogger

// init setup logger
func Setup() {
	// Define log level
	level := zap.NewAtomicLevel()
	level.SetLevel(zap.DebugLevel)

	// Create new zap logger
	logger = zap.New(zapcore.NewCore(
		getEncoder(),
		getLogWriter(),
		level,
	), zap.AddCaller())
	sugar = logger.Sugar()
	defer func() {
		logger.Sync()
		sugar.Sync()
	}()
}

func getEncoder() zapcore.Encoder {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	} // zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
	// return zapcore.NewConsoleEncoder(encoderCfg)
	return zapcore.NewJSONEncoder(encoderCfg)
}

func getLogWriter() zapcore.WriteSyncer {
	// file, _ := os.Create("./test.log")
	// return zapcore.AddSync(file)
	return zapcore.Lock(os.Stdout)
}

// Info log
func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Infof(template string, args ...interface{}) {
	sugar.Infof(template, args...)
}

// Debug log
func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func Debugf(template string, args ...interface{}) {
	sugar.Debugf(template, args...)
}

// Error log
func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func Errorf(template string, args ...interface{}) {
	sugar.Errorf(template, args...)
}

// Fatal log
func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}

func Fatalf(template string, args ...interface{}) {
	sugar.Fatalf(template, args...)
}
