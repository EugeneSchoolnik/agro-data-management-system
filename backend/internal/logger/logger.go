package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

// Init ініціалізує глобальний логер
func Init(level string) {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder // Читабельний час

	var encoder zapcore.Encoder
	if os.Getenv("APP_ENV") == "development" {
		encoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	} else {
		encoder = zapcore.NewJSONEncoder(config)
	}

	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zap.InfoLevel)
	Log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}