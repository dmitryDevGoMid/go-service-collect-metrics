package logger

import (
	"os"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

// Структура для реализации интерфейса
type ApiLogger struct {
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
func NewApiLogger(cfg *config.Config) *ApiLogger {
	return &ApiLogger{cfg: cfg}
}

// Дергаем уровень из карты, если не существует взвращаем zapcore.DebugLevel
func (l *ApiLogger) getLoggerLevel(cfg *config.Config) zapcore.Level {
	level, exist := loggerLevelMap[cfg.Logger.Level]
	if !exist {
		return zapcore.DebugLevel
	}

	return level
}

func (l *ApiLogger) InitLogger() {
	//Название файла куда пишем логи
	fileName := "log/zap.log"
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

func (l *ApiLogger) Debug(args ...interface{}) {
	l.sugarLogger.Debug(args...)
}

func (l *ApiLogger) Debugf(template string, args ...interface{}) {
	l.sugarLogger.Debugf(template, args...)
}

func (l *ApiLogger) Info(args ...interface{}) {
	l.sugarLogger.Info(args...)
}

func (l *ApiLogger) Infof(template string, args ...interface{}) {
	l.sugarLogger.Infof(template, args...)
}

func (l *ApiLogger) Printf(template string, args ...interface{}) {
	l.sugarLogger.Infof(template, args...)
}

func (l *ApiLogger) Warn(args ...interface{}) {
	l.sugarLogger.Warn(args...)
}

func (l *ApiLogger) Warnf(template string, args ...interface{}) {
	l.sugarLogger.Warnf(template, args...)
}

func (l *ApiLogger) Error(args ...interface{}) {
	l.sugarLogger.Error(args...)
}

func (l *ApiLogger) Errorf(template string, args ...interface{}) {
	l.sugarLogger.Errorf(template, args...)
}

func (l *ApiLogger) DPanic(args ...interface{}) {
	l.sugarLogger.DPanic(args...)
}

func (l *ApiLogger) DPanicf(template string, args ...interface{}) {
	l.sugarLogger.DPanicf(template, args...)
}

func (l *ApiLogger) Panic(args ...interface{}) {
	l.sugarLogger.Panic(args...)
}

func (l *ApiLogger) Panicf(template string, args ...interface{}) {
	l.sugarLogger.Panicf(template, args...)
}

func (l *ApiLogger) Fatal(args ...interface{}) {
	l.sugarLogger.Fatal(args...)
}

func (l *ApiLogger) Fatalf(template string, args ...interface{}) {
	l.sugarLogger.Fatalf(template, args...)
}
