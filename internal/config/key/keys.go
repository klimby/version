package key

// Viper keys.
const (
	AppName = "appName"
	Version = "version"
	WorkDir = "WORK_DIR"

	RemoteURL = "repoURL" // Remote repository URL.

	CfgFile = "configFile" // Configuration file name. Default: config.yaml.

	AllowCommitDirty      = "allowCommitDirty"      // Allow commit dirty repository. Default: false.
	AutoGenerateNextPatch = "autoGenerateNextPatch" // Auto generate next patch version, if version exists. Default: false.
	AllowDowngrades       = "allowDowngrades"       // Allow version downgrades. Default: false.

	GenerateChangelog   = "changelog.generate"   // Generate changelog. Default: true.
	ChangelogFileName   = "changelog.fileName"   // Changelog file name. Default: CHANGELOG.md.
	ChangelogTitle      = "changelog.title"      // Changelog title. Default: Changelog.
	ChangelogIssueURL   = "changelog.issueURL"   // Issue href template (with last slash). Default: empty.
	ChangelogShowAuthor = "changelog.showAuthor" // Show author in changelog. Default: false.
	ChangelogShowBody   = "changelog.showBody"   // Show body in changelog comment. Default: true.

	Backup  = "backupChanged" // Backup changed files. Default: false.
	Silent  = "silent"        // Silent mode from flags.
	DryRun  = "dryRun"        // Dry run mode from flags.
	Force   = "force"         // Force mode from flags.
	Verbose = "verbose"       // Verbose mode from flags.

	Prepare = "prepare" // Prepare flag in next command.
)

// Viper testing keys.
const (
	TestingSkipDIInit = "testing.skipDIInit" // Skip DI init. Default: false.
)
