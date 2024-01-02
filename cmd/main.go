package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/klimby/version/internal/backup"
	"github.com/klimby/version/internal/bump"
	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/internal/console"
	"github.com/klimby/version/internal/file"
	"github.com/klimby/version/internal/git"
	version2 "github.com/klimby/version/pkg/version"
	flags "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// -ldflags "-X main.version=$VERSION"
var version = "0.0.0"

func main() {
	ctx := context.Background()
	code := 0
	defer os.Exit(code)

	var (
		major       bool
		minor       bool
		patch       bool
		ver         string
		dryRun      bool
		silent      bool
		force       bool
		backupFiles bool
		test        bool

		removeBackup      bool
		generateConfig    bool
		generateChangelog bool
		configFile        string
	)

	f := file.NewFS()

	config.Init(func(options *config.ConfigOptions) {
		options.Version = version
	})

	repo, err := git.NewRepository(func(options *git.RepoOptions) {
		options.Path = viper.GetString(config.WorkDir)
	})
	if err != nil {
		console.Error(err.Error())
		code = 1

		return
	}

	remote, err := repo.RemoteURL()
	if err != nil || remote == "" {
		console.Warn("Remote repository URL is not set.")
	} else {
		viper.Set(config.RemoteURL, remote)
	}

	cfg, err := config.Load(f)
	if err != nil {
		fmt.Println(err)
		code = 1
		return
	}

	flagSet := flags.NewFlagSet("version", flags.ContinueOnError)
	flagSet.SortFlags = false

	flagSet.BoolVar(&major, "major", false, "Next major version")
	flagSet.BoolVar(&minor, "minor", false, "Next minor version")
	flagSet.BoolVar(&patch, "patch", false, "Next patch version")
	flagSet.BoolVar(&test, "test", false, "Next patch version")
	flagSet.StringVar(&ver, "version", "", "Version in format: 1.0.0")

	flagSet.BoolVarP(&dryRun, "dry", "d", false, "Dry run")
	flagSet.BoolVarP(&silent, "silent", "s", false, "Silent mode")
	flagSet.BoolVarP(&force, "force", "f", false, "Force mode")
	flagSet.BoolVarP(&backupFiles, "backup", "b", true, "Backup changed files")

	flagSet.BoolVar(&removeBackup, "remove-backup", false, "Remove backup files")
	flagSet.BoolVar(&generateConfig, "generate-config", false, "Generate or update config file")
	flagSet.BoolVar(&generateChangelog, "generate-changelog", false, "Generate and rewrite changelog file")
	flagSet.StringVarP(&configFile, "config", "c", "config.yaml", "Config file")

	flags.Parse()

	config.SetFlags(func(options *config.FlagOptions) {
		options.DryRun = dryRun
		options.Silent = silent
		options.Force = force
		options.Backup = backupFiles
	})

	if err := config.Check(cfg); err != nil {
		if errors.Is(err, config.ErrConfigWarn) {
			console.Warn(err.Error())
		} else {
			console.Error(err.Error())
			code = 1

			return
		}
	}

	switch {
	case generateConfig:
		fmt.Println("generate config")
		p := config.File(viper.GetString(config.ConfigFile))
		if err := backup.Create(f, p.Path()); err != nil {
			console.Error(err.Error())
			code = 1

			return
		}

		if err := config.Generate(f, cfg); err != nil {
			console.Error(err.Error())
			code = 1

			return
		}

		console.Success("Config file generated.")

	case generateChangelog:
		fmt.Println("generate changelog")

	case major:
		fmt.Println("major")

	case minor:
		fmt.Println("minor")

	case patch:
		fmt.Println("patch")

	case ver != "":
		fmt.Println("version")

	case removeBackup:
		fmt.Println("removeBackup")
		p := config.File(viper.GetString(config.ConfigFile))
		backup.Remove(f, p.Path())

		for _, bmp := range cfg.Bump {
			backup.Remove(f, bmp.File.Path())
		}

	case test:

		bumps := []config.BumpFile{
			{
				File: config.File("README.md"),
				RegExp: []string{
					`^!\[Version.*`,
				},
			},
		}

		ver := version2.V("1.0.0")

		viper.Set(config.Backup, true)

		bump.Apply(f, bumps, ver)

		fmt.Println("test")

	default:
		flagSet.PrintDefaults()
	}

	_ = ctx

	return

}
