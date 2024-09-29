package config

import (
	"encoding/json"
	"log/slog"
	"os"
)

func init() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: false})).With("group", "Config")
	if _, err := os.Stat("config.override.json"); os.IsNotExist(err) {
		log.Error("config.override.json not found")
	} else {
		data, err := os.ReadFile("config.override.json")
		if err != nil {
			log.Error("Read config.override.json", "error", err.Error())
			os.Exit(2)
		}

		err = json.Unmarshal(data, &Config)
		if err != nil {
			log.Error("Unmarshal config", "error", err.Error())
			os.Exit(2)
		}
	}
}
