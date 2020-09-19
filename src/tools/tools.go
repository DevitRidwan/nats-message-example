package tools

import (
	"flag"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func LoggerGenerator(OutputPaths []string, ErrorOutputPaths []string) (*zap.Logger, error) {
	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zapcore.InfoLevel),
		OutputPaths:      OutputPaths,
		ErrorOutputPaths: ErrorOutputPaths,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "info",
			LevelKey:     "level",
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			TimeKey:      "time",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	logger.Info("logger construction succeeded")
	return logger, err
}

func Ekstration(pathFlac string, nameFlac string, target ...interface{}) error {
	configPath := flag.String(pathFlac, ".", "Configuration YAML path")
	configName := flag.String(nameFlac, "config-prod", "Configuration Name { config-prod | config-dev } (Required)")
	flag.Parse()
	// config file path
	viper.AddConfigPath(*configPath)
	// config file name
	viper.SetConfigName(*configName)
	err := viper.ReadInConfig()

	if err != nil {
		return err
	}

	for _, element := range target {
		err = viper.Unmarshal(&element)
		if err != nil {
			return err
		}
	}
	return nil
}
