package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Path       string            `toml:"-"`
	Wallpapers string            `toml:"wallpapers"`
	Saved      string            `toml:"saved"`
	Current    string            `toml:"current"`
	Shell      string            `toml:"shell"`
	ShellArgs  []string          `toml:"shellArgs"`
	ClearEnv   bool              `toml:"clearEnv"`
	MaxHistory int               `toml:"maxHistory"`
	History    []string          `toml:"history"`
	Commands   map[string]string `toml:"commands"`
	Environ    map[string]string `toml:"environ"`
}

var defConfig = Config{
	Shell:      "bash",
	ShellArgs:  []string{"-c"},
	MaxHistory: 5,
	History:    []string{},
	Commands: map[string]string{
		"set":     "${basecmd} ${wpimg} --transition-type outer",
		"preview": "${basecmd} ${wpimg} --transition-type wipe --transition-angle 45",
		"restore": "${basecmd} ${wpimg} --transition-type wipe --transition-angle 225",
	},
	Environ: map[string]string{
		"basecmd": "awww img",
	},
}

func getConfigPath() string {
	if path := os.Getenv("WPMAN_CONFIG"); path != "" && filepath.IsAbs(path) {
		return path
	}

	cfg, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("Could not get user config dir:", err)
		os.Exit(1)
	}

	return filepath.Join(cfg, "wpmanrc.toml")
}

func loadConfig() Config {
	path := getConfigPath()
	if !fileExists(path) {
		fmt.Println("\x1b[2mConfig not found, creating new one with default values...\x1b[0m")

		writeConfig(defConfig, path)

		defConfig.Path = path
		return defConfig
	}

	var cfg Config
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Could not read config '"+path+"':", err)
		os.Exit(1)
	}

	if err := toml.Unmarshal(data, &cfg); err != nil {
		fmt.Println("Could not unmarshal config '"+path+"':", err)
		os.Exit(1)
	}

	cfg.Path = path
	return cfg
}

func writeConfig(cfg Config, path string) {
	data, err := toml.Marshal(cfg)
	if err != nil {
		fmt.Println("Could not marshal config:", err)
		os.Exit(1)
	}

	err = os.WriteFile(path, data, 0o644)
	if err != nil {
		fmt.Println("Could not write to config '"+path+"':", err)
		os.Exit(1)
	}
}
