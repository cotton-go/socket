package log

import (
	"context"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 定义常量 ctxLoggerKey,用于在上下文中存储 logger 实例
const ctxLoggerKey = "zapLogger"

// 定义 Logger 结构体，包含一个 *zap.Logger 类型的字段
type Logger struct {
	*zap.Logger
}

// Config 结构体用于存储配置信息
type Config struct {
	Env        string // 环境变量，如生产、开发等
	Level      string // 日志级别，如 debug、info、warn、error 等
	FileName   string // 日志文件名
	MaxSize    int    // 日志文件最大大小，单位为字节
	MaxBackups int    // 日志文件最大备份数
	MaxAge     int    // 日志文件最大保存时间，单位为天
	Compress   bool   // 是否启用日志文件压缩
	Encoding   string // 日志文件编码格式，如 json、console 等
}

// NewLog 创建一个新的日志记录器实例，根据配置文件中的设置进行初始化。
//
// 参数：
// - conf Config: 配置对象，包含日志记录器的配置信息。
//
// 返回值：
// - *Logger: 返回一个初始化完成的日志记录器实例。
func NewLog(conf Config) *Logger {
	// 根据配置文件中的日志级别设置日志级别
	var level zapcore.Level
	switch conf.Level {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}

	// 根据配置文件中的编码方式设置编码器
	var encoder zapcore.Encoder
	if conf.Encoding == "console" {
		encoder = zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "Logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseColorLevelEncoder,
			EncodeTime:     timeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.FullCallerEncoder,
		})
	} else {
		encoder = zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.EpochTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		})
	}

	// 将日志输出到控制台和文件中，如果指定了文件名，则还会写入到文件中。
	writer := zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)) // Print to console and file
	if conf.FileName != "" {
		hook := lumberjack.Logger{
			Filename:   conf.FileName,   // Log file path
			MaxSize:    conf.MaxSize,    // Maximum size unit for each log file: M
			MaxBackups: conf.MaxBackups, // The maximum number of backups that can be saved for log files
			MaxAge:     conf.MaxAge,     // Maximum number of days the file can be saved
			Compress:   conf.Compress,   // Compression or not
		}

		writer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)) // Print to console and file
	}

	core := zapcore.NewCore(encoder, writer, level)
	if conf.Env != "prod" {
		return &Logger{zap.New(core, zap.Development(), zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))}
	}

	return &Logger{zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))}
}

// timeEncoder 是一个时间编码器，它将给定的时间对象转换为特定格式的字符串并添加到编码器中。
//
// 参数：
//   - t time.Time: 需要编码的时间对象
//   - enc zapcore.PrimitiveArrayEncoder: 用于存储编码后的时间字符串的编码器
//
// 返回值：
//
//	无返回值，但会修改传入的编码器 enc
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	// 将时间对象转换为特定格式的字符串并添加到编码器中
	enc.AppendString(t.Format("2006-01-02 15:04:05.000000000"))
}

// WithValue 方法用于在给定的上下文中添加或更新字段值，并返回一个新的上下文。
//
// 参数：
//   - ctx context.Context: 需要添加或更新字段值的上下文
//   - fields ...zapcore.Field: 需要添加或更新的字段列表
//
// 返回值：
//   - context.Context: 包含添加或更新字段值的新上下文
func (l *Logger) WithValue(ctx context.Context, fields ...zapcore.Field) context.Context {
	// 使用 WithContext 方法获取当前上下文中的 logger,然后使用 With 方法添加或更新字段值
	return context.WithValue(ctx, ctxLoggerKey, l.WithContext(ctx).With(fields...))
}

// WithContext 方法用于从给定的上下文中获取 logger,并返回一个新的 logger。
//
// 参数：
//   - ctx context.Context: 需要获取 logger 的上下文
//
// 返回值：
//   - *Logger: 从上下文中获取的 logger
func (l *Logger) WithContext(ctx context.Context) *Logger {
	// 从上下文中获取 logger
	zl := ctx.Value(ctxLoggerKey)
	// 如果获取到的 logger 是 zap.Logger 类型，则创建一个新的 Logger 并返回
	if ctxLogger, ok := zl.(*zap.Logger); ok {
		return &Logger{ctxLogger}
	}
	// 否则直接返回当前的 logger
	return l
}
