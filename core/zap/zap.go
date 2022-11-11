package zap

import (
	"fmt"
	"github.com/hihibug/microdule/core/utils"
	zaprotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"time"
)

type (
	Zap struct {
		Log   *zap.Logger
		Conf  *Config
		level zapcore.Level
	}

	Log interface {
		Client() *zap.Logger
		GetWriteSyncer() (zapcore.WriteSyncer, error)
		GetEncoderCore() (core zapcore.Core)
		GetEncoderConfig() (config zapcore.EncoderConfig)
		GetEncoder() zapcore.Encoder
		CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder)
	}
)

func NewZap(conf *Config) Log {

	err := conf.Validate()
	if err != nil {
		panic(err)
	}

	if ok, _ := utils.PathExists(conf.Director); !ok { // 判断是否有Director文件夹
		_ = os.Mkdir(conf.Director, os.ModePerm)
	}
	var z Zap

	z.Conf = conf

	switch conf.Level { // 初始化配置文件的Level
	case "debug":
		z.level = zap.DebugLevel
	case "info":
		z.level = zap.InfoLevel
	case "warn":
		z.level = zap.WarnLevel
	case "error":
		z.level = zap.ErrorLevel
	case "dpanic":
		z.level = zap.DPanicLevel
	case "panic":
		z.level = zap.PanicLevel
	case "fatal":
		z.level = zap.FatalLevel
	default:
		z.level = zap.InfoLevel
	}

	if z.level == zap.DebugLevel || z.level == zap.ErrorLevel {
		z.Log = zap.New(z.GetEncoderCore(), zap.AddStacktrace(z.level))
	} else {
		z.Log = zap.New(z.GetEncoderCore())
	}

	if conf.ShowLine {
		z.Log = z.Log.WithOptions(zap.AddCaller())
	}

	return &z
}

func (z *Zap) Client() *zap.Logger {
	return z.Log
}

// GetWriteSyncer 使用file-rotatelogs进行日志分割
func (z *Zap) GetWriteSyncer() (zapcore.WriteSyncer, error) {
	fileWriter, err := zaprotatelogs.New(
		path.Join(z.Conf.Director, "%Y-%m-%d.log"),
		zaprotatelogs.WithMaxAge(7*24*time.Hour),
		zaprotatelogs.WithRotationTime(24*time.Hour),
	)
	if z.Conf.LogInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter)), err
	}
	return zapcore.AddSync(fileWriter), err
}

// GetEncoderCore 获取 Encoder 的 zapcore.Core
func (z *Zap) GetEncoderCore() (core zapcore.Core) {
	writer, err := z.GetWriteSyncer()
	if err != nil {
		panic("Get Write Syncer Failed err:" + err.Error())
	}

	return zapcore.NewCore(z.GetEncoder(), writer, z.level)
}

// GetEncoderConfig 获取 zapcore.EncoderConfig
func (z *Zap) GetEncoderConfig() (config zapcore.EncoderConfig) {
	config = zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  z.Conf.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     z.CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	switch {
	case z.Conf.EncodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	case z.Conf.EncodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		config.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case z.Conf.EncodeLevel == "CapitalLevelEncoder": // 大写编码器
		config.EncodeLevel = zapcore.CapitalLevelEncoder
	case z.Conf.EncodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	return config
}

// GetEncoder 获取 zapcore.Encoder
func (z *Zap) GetEncoder() zapcore.Encoder {
	if z.Conf.Format == "json" {
		return zapcore.NewJSONEncoder(z.GetEncoderConfig())
	}
	return zapcore.NewConsoleEncoder(z.GetEncoderConfig())
}

// CustomTimeEncoder 自定义日志输出时间格式
func (z *Zap) CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(z.Conf.Prefix + " 2006/01/02 - 15:04:05.000"))
}

func NewZapWriter(z *zap.Logger) *Zap {
	return &Zap{
		Log: z,
	}
}

func (z *Zap) Printf(format string, v ...interface{}) {
	z.Log.Warn(fmt.Sprintf(format+" ", v...))
	return
}
