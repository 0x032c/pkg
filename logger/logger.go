package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	zapLogger *zap.Logger
)

type Config struct {
	LogPath    string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Level      string
}

func DefaultConfig() Config {
	return Config{
		LogPath:    "./logs/app.log",
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     30,
		Level:      "info",
	}
}

func InitLogger(cfg Config) error {
	if cfg.LogPath == "" {
		return fmt.Errorf("log path is required")
	}
	if err := os.MkdirAll(filepath.Dir(cfg.LogPath), 0755); err != nil {
		return err
	}
	fileSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.LogPath,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
	})
	consoleSyncer := zapcore.AddSync(os.Stdout)
	encCfg := zapcore.EncoderConfig{
		MessageKey: "msg", LevelKey: "level", TimeKey: "ts", CallerKey: "caller",
		EncodeLevel: zapcore.CapitalLevelEncoder, EncodeTime: zapcore.ISO8601TimeEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder, LineEnding: zapcore.DefaultLineEnding,
	}
	fileEncoder := zapcore.NewJSONEncoder(encCfg)
	consoleEncoderCfg := encCfg
	consoleEncoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderCfg)
	level := zapcore.InfoLevel
	_ = level.UnmarshalText([]byte(strings.ToLower(cfg.Level)))
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, fileSyncer, level),
		zapcore.NewCore(consoleEncoder, consoleSyncer, level),
	)
	zapLogger = zap.New(core, zap.AddCaller())
	return nil
}

func Logger() *zap.Logger {
	if zapLogger == nil {
		log.Println("Logger not initialized")
		return zap.NewNop()
	}
	return zapLogger
}

func Sync() error {
	if zapLogger != nil {
		return zapLogger.Sync()
	}
	return nil
}

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		Logger().Info("HTTP",
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("ip", c.ClientIP()),
			zap.String("ua", c.Request.UserAgent()),
			zap.Duration("latency", time.Since(start)),
		)
	}
}

func GinRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				Logger().Error("Panic", zap.Any("err", err))
				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}

func Info(msg string, fields ...zap.Field)  { Logger().Info(msg, fields...) }
func Error(msg string, fields ...zap.Field) { Logger().Error(msg, fields...) }
func Debug(msg string, fields ...zap.Field) { Logger().Debug(msg, fields...) }
func Warn(msg string, fields ...zap.Field)  { Logger().Warn(msg, fields...) }
func Fatal(msg string, fields ...zap.Field) { Logger().Fatal(msg, fields...) }

func Infof(format string, args ...interface{})  { Logger().Info(fmt.Sprintf(format, args...)) }
func Errorf(format string, args ...interface{}) { Logger().Error(fmt.Sprintf(format, args...)) }
func Debugf(format string, args ...interface{}) { Logger().Debug(fmt.Sprintf(format, args...)) }
func Warnf(format string, args ...interface{})  { Logger().Warn(fmt.Sprintf(format, args...)) }
func Fatalf(format string, args ...interface{}) { Logger().Fatal(fmt.Sprintf(format, args...)) }
