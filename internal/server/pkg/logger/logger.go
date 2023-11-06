package logger

import (
	"os"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/config"
	"gopkg.in/natefinch/lumberjack.v2"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Структура для реализации интерфейса
type APILogger struct {
	cfg         *config.Config
	sugarLogger *zap.SugaredLogger
}

// Методы логирования
type Logger interface {
	InitLogger()
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	DPanic(args ...interface{})
	DPanicf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
	Printf(template string, args ...interface{})
}

// Перемеррам с уровнями логирования
var loggerLevelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

// Конструктор для логгера возвращает структура с методами интерфейса
func NewAPILogger(cfg *config.Config) *APILogger {
	return &APILogger{cfg: cfg}
}

// Дергаем уровень из карты, если не существует взвращаем zapcore.DebugLevel
func (l *APILogger) getLoggerLevel(cfg *config.Config) zapcore.Level {
	level, exist := loggerLevelMap[cfg.Logger.Level]
	if !exist {
		return zapcore.DebugLevel
	}

	return level
}

func (l *APILogger) InitLogger() {
	//Название файла куда пишем логи
	fileName := "./internal/server/log/zap.log"
	//Получаем уровень логирования
	logLevel := l.getLoggerLevel(l.cfg)

	//logWriter := zapcore.AddSync(os.Stderr)
	//Выполняем ротацию журналов с помощью lumberjack
	//Пишем логи в файл

	logWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename: fileName,
		//MaxSize:   1 << 30, //1G
		MaxSize:   1, //1M
		LocalTime: true,
		Compress:  true,
	})

	//var encoderCfg zapcore.EncoderConfig

	encoderCfg := zap.NewDevelopmentEncoderConfig()

	//EncoderConfig устанавливаем параметры для кодировщика
	encoder := make(map[string]zapcore.Encoder)

	//Ключи
	encoderCfg.LevelKey = "LEVEL"
	encoderCfg.CallerKey = "CALLER"
	encoderCfg.TimeKey = "TIME"
	encoderCfg.NameKey = "NAME"
	encoderCfg.MessageKey = "MESSAGE"

	//Определяем тип логирования
	encoder["console"] = zapcore.NewConsoleEncoder(encoderCfg)
	encoder["json"] = zapcore.NewJSONEncoder(encoderCfg)

	//Параметры для кодировщика формат времени
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	var core zapcore.Core

	//NewCore создает ядро, которое записывает журналы в WriteSyncer.
	if l.cfg.Logger.Encoding == "full" {
		//Пишем как в консоль так и в файл
		core = zapcore.NewTee(zapcore.NewCore(encoder["json"], logWriter, zap.NewAtomicLevelAt(logLevel)),
			zapcore.NewCore(encoder["console"], zapcore.AddSync(os.Stdout), zap.NewAtomicLevelAt(logLevel)))
	} else if l.cfg.Logger.Encoding == "console" {
		// Пишем только в консоль
		core = zapcore.NewCore(encoder["console"], zapcore.AddSync(os.Stdout), zap.NewAtomicLevelAt(logLevel))
	} else {
		// Пишем только в файл
		core = zapcore.NewCore(encoder["json"], logWriter, zap.NewAtomicLevelAt(logLevel))
	}

	//core := zapcore.NewCore(encoder["console"], zapcore.AddSync(os.Stdout), zap.NewAtomicLevelAt(logLevel))

	//New создает новый регистратор из предоставленных zapcore.Core и Options.
	//Если переданный zapcore.Core равен нулю, он возвращается к использованию бездействующей реализации.
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	//Синтаксичекий сахар для удобного ведения журналов
	l.sugarLogger = logger.Sugar()

	//Запускаем сахарный логер
	if err := l.sugarLogger.Sync(); err != nil {
		l.sugarLogger.Error(err)
	}
}

// Методы логирования

func (l *APILogger) Debug(args ...interface{}) {
	l.sugarLogger.Debug(args...)
}

func (l *APILogger) Debugf(template string, args ...interface{}) {
	l.sugarLogger.Debugf(template, args...)
}

func (l *APILogger) Info(args ...interface{}) {
	l.sugarLogger.Info(args...)
}

func (l *APILogger) Infof(template string, args ...interface{}) {
	l.sugarLogger.Infof(template, args...)
}

func (l *APILogger) Printf(template string, args ...interface{}) {
	l.sugarLogger.Infof(template, args...)
}

func (l *APILogger) Warn(args ...interface{}) {
	l.sugarLogger.Warn(args...)
}

func (l *APILogger) Warnf(template string, args ...interface{}) {
	l.sugarLogger.Warnf(template, args...)
}

func (l *APILogger) Error(args ...interface{}) {
	l.sugarLogger.Error(args...)
}

func (l *APILogger) Errorf(template string, args ...interface{}) {
	l.sugarLogger.Errorf(template, args...)
}

func (l *APILogger) DPanic(args ...interface{}) {
	l.sugarLogger.DPanic(args...)
}

func (l *APILogger) DPanicf(template string, args ...interface{}) {
	l.sugarLogger.DPanicf(template, args...)
}

func (l *APILogger) Panic(args ...interface{}) {
	l.sugarLogger.Panic(args...)
}

func (l *APILogger) Panicf(template string, args ...interface{}) {
	l.sugarLogger.Panicf(template, args...)
}

func (l *APILogger) Fatal(args ...interface{}) {
	l.sugarLogger.Fatal(args...)
}

func (l *APILogger) Fatalf(template string, args ...interface{}) {
	l.sugarLogger.Fatalf(template, args...)
}
