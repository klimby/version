package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/klimby/version/internal/action"
	"github.com/klimby/version/internal/bump"
	"github.com/klimby/version/internal/changelog"
	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/internal/console"
	"github.com/klimby/version/internal/file"
	"github.com/klimby/version/internal/git"
	vers "github.com/klimby/version/pkg/version"
	flags "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// -ldflags "-X main.version=$VERSION"
var version = "0.0.0"

func main() {
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
		code = errorHandler(err)

		return
	}

	remote, err := repo.RemoteURL()
	if err != nil {
		console.Warn(err.Error())
	}

	config.SetUrlFromGit(remote)

	cfg, err := config.Load(f)
	if err != nil {
		code = errorHandler(err)

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
	flagSet.BoolVarP(&backupFiles, "backup", "b", false, "Backup changed files")

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
			code = errorHandler(err)

			return
		}
	}

	chGen := changelog.NewGenerator(f, repo)

	if viper.GetBool(config.DryRun) {
		console.Info("Dry run mode.")
	}

	actions := action.New(func(o *action.Args) {
		o.Repo = repo
		o.Config = cfg
		o.ChGen = chGen
		o.ReadWriter = f
	})

	switch {
	case generateConfig:
		console.Notice("Generate config file\n")

		if err := actions.GenerateConfig(); err != nil {
			code = errorHandler(err)

			return
		}

		console.Success("Config file generated.")

	case generateChangelog:
		console.Notice("Generate changelog\n")
		if err := actions.GenerateChangelog(); err != nil {
			code = errorHandler(err)

			return
		}

		console.Success("Changelog generated.")

	case major:
		console.Notice("Bump major version\n")
		if err := actions.NextVersion(git.NextMajor, ""); err != nil {
			code = errorHandler(err)

			return
		}

		console.Success("Version bumped.")

	case minor:
		console.Notice("Bump minor version\n")

		if err := actions.NextVersion(git.NextMinor, ""); err != nil {
			code = errorHandler(err)

			return
		}

		console.Success("Version bumped.")

	case patch:
		console.Notice("Bump patch version\n")

		if err := actions.NextVersion(git.NextPatch, ""); err != nil {
			code = errorHandler(err)

			return
		}

		console.Success("Version bumped.")

	case ver != "":
		console.Notice("Bump version to " + ver + "\n")

		if err := actions.NextVersion(git.NextCustom, vers.V(ver)); err != nil {
			code = errorHandler(err)

			return
		}

		console.Success("Version bumped.")

	case removeBackup:
		console.Notice("Remove backup files\n")
		actions.RemoveBackup(f)
		console.Success("Backup files removed.")

	case test:

		bumps := []config.BumpFile{
			{
				File: config.File("README.md"),
				RegExp: []string{
					`^!\[Version.*`,
				},
			},
		}

		ver := vers.V("1.0.0")

		viper.Set(config.Backup, true)

		bump.Apply(f, bumps, ver)

		fmt.Println("test")

	default:
		flagSet.PrintDefaults()
	}

	return

}

func errorHandler(err error) int {
	if err != nil {
		console.Error(err.Error())

		return 1
	}

	return 0
}
