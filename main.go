package main

import (
	"fmt"
	"maps"
	"math/rand/v2"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
)

var helpMessage = `
WallPaper MANager

Usage:
  wpman [COMMAND] [ARGS] [extra]

Commands:
  set [image] [extra]     sets image as wallpaper and saves the path *
  save [image?]           saves image path *
  preview [image] [extra] temporarily sets wallpaper to provided path *
  restore                 changes wallpaper to saved path
  get                     outputs current wallpaper path

* if image parameter equals 'random', then random wallpaper from wallpaper dir is applied.
`

func main() {
	cfg := loadConfig()

	if len(os.Args) < 2 {
		fmt.Println(helpMessage)
		os.Exit(2)
	}

	if !checkConfig(cfg, os.Args[1]) {
		os.Exit(78)
	}

	varRegex = regexp.MustCompile(`\$[0-9_a-zA-Z]`)

	switch os.Args[1] {
	case "set":
		if len(os.Args) < 3 {
			fmt.Println("You should provide image path or 'random' as second argument.")
			os.Exit(2)
		}

		path := os.Args[2]
		if path == "random" {
			path = randomWp(&cfg)
		}

		run(cfg, cfg.Commands["set"], path, os.Args[3:])
		cfg.Current = path
		cfg.Saved = path

	case "save":
		var path string
		if len(os.Args) > 2 {
			path = os.Args[2]

			if path == "random" {
				path = randomWp(&cfg)
			}
		} else if cfg.Current != "" {
			path = cfg.Current
		} else {
			fmt.Println("There is no current image and no image was provided by CLI.")
			os.Exit(1)
		}

		if !fileExists(path) {
			fmt.Println("Provided image (" + path + ") doesn't exist.")
			os.Exit(66)
		}
		cfg.Saved = path
		fmt.Println("Saved '" + path + "'")

	case "preview":
		if len(os.Args) < 3 {
			fmt.Println("You should provide image path or 'random' as second argument.")
			os.Exit(2)
		}

		path := os.Args[2]
		if path == "random" {
			path = randomWp(&cfg)
		}

		run(cfg, cfg.Commands["preview"], path, os.Args[3:])

		fmt.Println("Previewing '" + path + "'")
		cfg.Current = path

	case "restore":
		run(cfg, cfg.Commands["restore"], cfg.Saved, os.Args[2:])

		fmt.Println("Restored to '" + cfg.Saved + "'")
		cfg.Current = cfg.Saved

	case "get":
		fmt.Print(cfg.Current)

	default:
		fmt.Println(helpMessage)
		os.Exit(2)
	}

	writeConfig(cfg, cfg.Path)
}

func run(cfg Config, cmd string, arg string, vars []string) {
	env := maps.Clone(cfg.Environ)
	if env == nil {
		env = make(map[string]string)
	}
	env["wpimg"] = arg

	// inject args to vars
	for i, v := range vars {
		env["arg"+strconv.Itoa(i)] = v
	}

	exe := expandVars(cmd, env, cfg.ClearEnv)

	if output, err := exec.Command(cfg.Shell, append(cfg.ShellArgs, exe)...).Output(); err != nil {
		fmt.Println("Error while executing command:\n  "+cmd+"\nOutput:\n  "+string(output)+"\nError:\n ", err)
		os.Exit(1)
	}
}

func randomWp(cfg *Config) string {
	// clear the slice so history contains only cfg.MaxHistory paths
	for len(cfg.History) != 0 && len(cfg.History) >= cfg.MaxHistory {
		cfg.History = slices.Delete(cfg.History, 0, 1)
	}

	dir, err := os.ReadDir(cfg.Wallpapers)
	if err != nil {
		fmt.Println("Could not open wallpapers dir:", err)
		os.Exit(1)
	}

	candidates := []string{}
	for _, e := range dir {
		if !e.IsDir() {
			candidates = append(candidates, e.Name())
		}
	}

	for i := 0; i < len(candidates); i++ {
		r := rand.IntN(len(candidates))
		if !slices.Contains(cfg.History, candidates[r]) {
			cfg.History = append(cfg.History, candidates[r])
			return filepath.Join(cfg.Wallpapers, candidates[r])
		}
		candidates = slices.Delete(candidates, r, r)
	}

	fmt.Println("Could not find unique candidates for random. Every candidate already exists in history.")
	os.Exit(1)
	return ""
}

func checkConfig(cfg Config, cmd string) bool {
	problem := true
	if cfg.Wallpapers == "" && (cmd == "set" || cmd == "preview" || cmd == "save") {
		fmt.Println("You need to specify wallpapers path to use wpman " + cmd + ".\nPlease, fill out 'wallpapers' field in your config")
	} else if (cmd == "set" || cmd == "preview" || cmd == "save") && !fileExists(cfg.Wallpapers) {
		fmt.Println("Provided wallpapers path does not exist (" + cfg.Wallpapers + ")")
	} else if cfg.Shell == "" && (cmd == "set" || cmd == "preview" || cmd == "restore") {
		fmt.Println("You should specify shell to use. Also, it'd be a good idea to add -c to 'shell_args'")
	} else if s, ok := cfg.Commands["set"]; cmd == "set" && (!ok || s == "") {
		fmt.Println("'set' command should be specified in config.")
	} else if s, ok := cfg.Commands["preview"]; cmd == "preview" && (!ok || s == "") {
		fmt.Println("'preview' command should be specified in config.")
	} else if s, ok := cfg.Commands["restore"]; cmd == "restore" && (!ok || s == "") {
		fmt.Println("'restore' command should be specified in config.")
	} else {
		problem = false
	}

	if problem {
		fmt.Println("Config path: '" + cfg.Path + "'")
		return false
	}
	return true
}
