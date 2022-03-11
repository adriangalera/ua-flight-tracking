package config

import (
	"agalera.eu/flight-tracking/internal/zone"
	"github.com/spf13/viper"
	"log"
)

type Configuration struct {
	Zones []zone.Zone
}

func ReadConfiguration(filename string) (*Configuration, error) {
	viper.SetConfigFile(filename)
	var newConfig Configuration

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file, %s", err)
		return nil, err
	}

	_ = viper.Unmarshal(&newConfig)
	return &newConfig, nil
}
