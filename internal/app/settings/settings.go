package settings

import (
	"io"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v3"
)

var paths = [2]string{"./app.yml", "./configs/app.yml"}

type logs struct {
	FileEnabled bool `yaml:"file-enabled"`
	MaxSize     int  `yaml:"max-size"`
	MaxAge      int  `yaml:"max-age"`
}

type server struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type database struct {
	File string `yaml:"file"`
}

type Settings struct {
	Logs     logs     `yaml:"logs"`
	Server   server   `yaml:"server"`
	Database database `yaml:"database"`
}

var loadedSettings *Settings

func defaultSettings() *Settings {
	return &Settings{
		Logs: logs{
			FileEnabled: true,
			MaxSize:     10,
			MaxAge:      30,
		},
		Server: server{
			Host: "0.0.0.0",
			Port: "8080",
		},
		Database: database{
			File: "./data/app.db",
		},
	}
}

func (s *Settings) configureLogger() {
	var handler *slog.TextHandler
	options := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	if s.Logs.FileEnabled {
		logsFile := &lumberjack.Logger{
			Filename:  "./logs/app.log",
			MaxSize:   s.Logs.MaxSize,
			MaxAge:    s.Logs.MaxAge,
			LocalTime: true,
		}
		multiWriter := io.MultiWriter(os.Stdout, logsFile)
		handler = slog.NewTextHandler(multiWriter, options)
	} else {
		handler = slog.NewTextHandler(os.Stdout, options)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func LoadSettings() Settings {
	if loadedSettings != nil {
		slog.Error("Settings already loaded")
		os.Exit(1)
	}

	loadedSettings = defaultSettings()
	var file *os.File
	var fileErr error

	for _, path := range paths {
		file, fileErr = os.Open(path)

		if fileErr != nil {
			slog.Warn("Settings not found", "path", path)
		} else {
			slog.Info("Settings found", "path", path)
			break
		}
	}

	if file != nil {
		defer file.Close()
		decoder := yaml.NewDecoder(file)

		if err := decoder.Decode(&loadedSettings); err != nil {
			slog.Error("Decoding YAML settings", "error", err)
			os.Exit(1)
		}
	}

	loadedSettings.configureLogger()

	return *loadedSettings
}
