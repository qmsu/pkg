package log

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DebugLevel int8 = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
	DPanicLevel
	PanicLevel
	FatalLevel
)

var log *zap.Logger

func Sugar() *zap.SugaredLogger {
	return log.Sugar()
}

type Config struct {
	Level      int8   `toml:"level"`
	OutputPath string `toml:"output_path"` //日志文件路径
	MaxSize    int    `toml:"max_size"`    // 每个文件大小 M
	MaxBackups int    `toml:"max_backups"` // 最多保留30个备份
	MaxAge     int    `toml:"max_age"`     // 日志文件保存时间 天
	Compress   bool   `toml:"compress"`    // 是否压缩 disabled by default
}

func Init(conf Config) {
	if log == nil {
		lumberJackLogger := &lumberjack.Logger{
			Filename:   conf.OutputPath, // 日志文件路径
			MaxSize:    conf.MaxSize,    // 每个文件最大字节
			MaxBackups: conf.MaxBackups, // 最多保留30个备份
			MaxAge:     conf.MaxAge,     // 日志文件保存时间
			Compress:   conf.Compress,   // 是否压缩 disabled by default
		}
		writeSyncer := zapcore.AddSync(lumberJackLogger)
		encoder := getEncoder()
		core := zapcore.NewCore(encoder, writeSyncer, getLevel(conf.Level))
		log = zap.New(core, zap.AddCaller(), zap.Development())
		zap.ReplaceGlobals(log) // 替换zap包中全局的logger实例，后续在其他包中只需使用zap.L()调用即可
	}
	return
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		log.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					log.Sugar().Errorw(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					log.Sugar().Errorw("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					log.Sugar().Errorw("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

func getLevel(level int8) zapcore.Level {
	switch level {
	case DebugLevel:
		return zap.DebugLevel
	case InfoLevel:
		return zap.InfoLevel
	case WarnLevel:
		return zap.WarnLevel
	case ErrorLevel:
		return zap.ErrorLevel
	default:
		return zap.ErrorLevel
	}
}
