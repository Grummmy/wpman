# WallPaper MANager

Most linux distros have wallpapers. Hyprland users might prefer [awww](https://codeberg.org/LGFae/awww), Gnome users have built in app for wallpapers, as well as KDE Plasma users. Wallpaper apps can set wallpapers, maybe sat them with animation (like awww) or set videos as wallpapers ([mpvpaper](https://github.com/GhostNaN/mpvpaper)). But, most of wallpaper apps, or every that I know, don't manage them further then just set it as background layer.

If you want to randomise your wallpaper - you write a custom script. If you want to preview the wallpaper in use and then revert back - you write a custom script. If you want to get current wallpaper image - you write a custom wrapper script around your wallpaper app. And I created a (almost) universal solution - WPMAN.

#### Can I use it?
TL;DR: Most likely yes if you use Linux or MacOS. You can use it on Windows if you find a way to set wallpaper with a command.

Basically, the program is cross-platform, since it's written in Go. But, in order to work WPMAN runs shell commands to change wallpaper. **Real requirements** are a shell (ex. bash) and a way to change wallpaper with a command.

## Features
- Set wallpaper and save it
- Preview questionable wallpaper, don't save it
- Restore saved wallpaper
- Save wallpaper path
- Get wallpaper path
- Any wallpaper backend

    WPMAN itself doesn't do anything with wallpapers`*`, it only runs preconfigured commands with arguments. If you wallpaper backend can be used via CLI - it works.

- Customizable commands

    You can pass extra parameters to cli, image path isn't the only dynamic part of a command. WPMAN supports 3 different commands to set, preview and restore wallpaper so you can set up commands with different animations.

- Environment variables control

    Run commands with clear env vars or add your own ones in config.

- Random wallpaper

    Instead of specifying direct image path you can type `random`, and WPMAN will pick random image from your configured wallpapers directory.

- Random history

    When you use `random` instead of image path, final image is added to history. When picking next random image WPMAN makes sure it's not a recent one. History has max size, when reached first added image will be deleted from history.

`*` WPMAN only runs predefined commands in config. You should already have working wallpaper backend in order to use WPMAN.

## Instalation

Clone the repo and cd in it.
```shell
git clone https://github.com/Grummmy/wpman.git
cd wpman
```

Build it, you need [Go](https://go.dev/) installed. `ldflags` and `trimpath` are used to reduce the size of resulting binary.
```shell
go build -o wpman -trimpath -ldflags "-w -s" .
```

Copy the binary to your PATH and you are good to go.
```shell
cp wpman ~/.local/bin/
# or, if you want to install it system-wide
sudo cp wpman /usr/local/bin/
```


## Usage

```
wpman [COMMAND] [ARGS] [EXTRA]

Commands
  set [image] [extra]     sets image as wallpaper and saves the path *
  save [image?]           saves image path *
  preview [image] [extra] temporarily sets wallpaper to provided path *
  restore                 changes wallpaper to saved path
  get                     outputs current wallpaper path

* if image parameter equals 'random', then random wallpaper from wallpaper dir is applied.
```

**⚠ Note**: `get` only outputs correct current wallpaper if wallpaper was set by WPMAN. If you set wallpaper via awww, gnome settings or any other app - it will desync and willoutput incorrect value. Use WPMAN to set (or preview) wallpaper and it will fix the problem.


## Configuration

Config is written in [toml](https://toml.io/). Config is overwritten every WPMAN run, comments are not preserved.

Default config location is `$XDG_CONFIG_HOME/wpmanrc.toml` (often is `~/.config/wpmanrc.toml`). If `WPMAN_CONFIG` environment variable is set to a valid filepath and file exists, it is used instead of default location.

Example configuration
```toml
wallpapers = ""
saved = ""
current = ""
shell = "bash"
shellArgs = ["-c"]
clearEnv = false
maxHistory = 5
history = []

[commands]
  preview = "${basecmd} ${wpimg} --transition-type wipe --transition-angle 45"
  restore = "${basecmd} ${wpimg} --transition-type wipe --transition-angle 225"
  set = "${basecmd} ${wpimg} --transition-type wipe --transition-angle 225"

[environ]
  basecmd = "awww img"
```

### wallpapers
Path to folder that contains wallpaper images. `random` picks image from this folder.
Should be absolute path, for example `/home/user/Pictures`.

### saved
Path to saved wallpaper. Used by WPMAN, not designed to be changed by user.

### current
Path to current wallpaper. Commands set, preview, and restore change it. Used by WPMAN, not designed to be changed by user.

### shell
Shell to use when running commands. Can be any valid shell.

### shellArgs
Argument to pass when running shell before passing the command. `-c` is added by default because in most shells (bash, zsh, fish) it executes passed command.

### clearEnv
Whether to clear environment variables before running wallpaper command through shell. This parameter doesn't affect custom environment variabels specified in config and extra variables passed in CLI.

### maxHistory
Maximum random history entries to save. Set to 0 or below to disable history. Config will still contain last set item in history, but it won't be used by WPMAN.

### history
List of recent files from wallpapers directory selected by `random`. See previous value to set maximum capacity. Used by WPMAN, not designed to be changed by user.

### commands
Contains commands to set, preview and restore wallpaper. Commands should be non-empty when calling these commands from WPMAN. The command is passed as single argument after all shell arguments.

Commands can have variables in them. Variables are expanded with this `${var}` syntax, where `var` is variable name. Variables are expanded by `os.ExpandEnv` function, but `$var` syntax is disabled so you can use dollar sign (`$`) with letters. Here is a list of variables you can use:

- `wpimg` var is replaced with image path being used in command.
- `argN` var is replaced with extra arguments that your add in CLI after image path. `N` is replaced with argument number starting from 0. If you run `wpman set random first anotherOneArg`, then `arg0` will be `first` and `arg1` will be `anotherOneArg`.
- You can use environment variables, for example `XDG_CONFIG_HOME` or `HOME`. Unless `clearEnv` is set to `true`.
- Environment variables declared in config.

### environ
Contains extra environment variables that get injected before expanding variables. Variables here are not affected by `clearEnv`.
