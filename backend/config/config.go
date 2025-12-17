package config

import (
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// testing
type Source struct {
	ID      string `yaml:"-" json:"id"`
	Name    string `yaml:"name" json:"name"`
	Path    string `yaml:"path" json:"path"`
	Type    string `yaml:"-" json:"type"`
	Enabled bool   `yaml:"-" json:"enabled"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type Config struct {
	Server  ServerConfig `yaml:"server"`
	Sources []Source     `yaml:"sources"`
}

var AppConfig Config

func Init() {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal("Failed to read config.yaml:", err)
	}

	if err := yaml.Unmarshal(data, &AppConfig); err != nil {
		log.Fatal("Failed to parse config.yaml:", err)
	}

	// Auto-generate IDs and set defaults
	for i := range AppConfig.Sources {
		AppConfig.Sources[i].ID = generateID(AppConfig.Sources[i].Name)
		AppConfig.Sources[i].Type = "local"
		AppConfig.Sources[i].Enabled = true
	}

	log.Printf("Loaded %d sources from config", len(AppConfig.Sources))
}

func GetEnabledSources() []Source {
	var enabled []Source
	for _, src := range AppConfig.Sources {
		if src.Enabled {
			enabled = append(enabled, src)
		}
	}
	return enabled
}

func generateID(name string) string {
	// Convert "Google Drive" -> "google-drive"
	id := strings.ToLower(name)
	id = strings.ReplaceAll(id, " ", "-")
	return id
}
