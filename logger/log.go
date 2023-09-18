package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Logger      *zap.SugaredLogger
	writeSyncer zapcore.WriteSyncer
	encoder     zapcore.Encoder
	core        zapcore.Core
	logger      *zap.Logger
)

const defaultLogPath = "guardian.log"

func init() {
	writeSyncer = getLogWriter(defaultLogPath)
	encoder = getEncoder()
	core = zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel)
	logger = zap.New(core, zap.AddCaller())
	Logger = logger.Sugar()
	_ = Logger.Sync()
}

func InitLogger(maxSize, maxRolls int, path, level string) {
	writeSyncer = getLogWriter(path)
	encoder = getEncoder()
	core = zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel)
	lumberjackLogger := &lumberjack.Logger{
		Filename:   path,
		MaxSize:    maxSize,
		MaxBackups: maxRolls,
	}
	writeSyncer = zapcore.AddSync(lumberjackLogger)
	logLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		logLevel = zapcore.InfoLevel
	}
	core = zapcore.NewCore(encoder, writeSyncer, logLevel)
	logger = zap.New(core, zap.AddCaller())
	Logger = logger.Sugar()
	_ = Logger.Sync()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(path string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   path,
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}
