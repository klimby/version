package config

// Viper keys.
const (
	appName = "appName"
	Version = "version"
	WorkDir = "WORK_DIR"

	RemoteURL = "repoURL" // Remote repository URL.

	RunBefore = "runBefore" // Run command before commit. []string. Default: empty.
	RunAfter  = "runAfter"  // Run command after commit. []string. Default: empty.

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

	changelogCommitTypes = "changelog.commitTypes" // Commit types for changelog. Map[string]string. Key - type key, value- type name.
	changelogCommitOrder = "changelog.commitOrder" // Commit types order for changelog. []string. Default: empty.

	Backup  = "backupChanged" // Backup changed files. Default: false.
	Silent  = "silent"        // Silent mode from flags.
	DryRun  = "dryRun"        // Dry run mode from flags.
	Force   = "force"         // Force mode from flags.
	Verbose = "verbose"       // Verbose mode from flags.
)

// Default values.
const (
	_AppName               = "Version"
	_Version               = "0.0.0"
	_AllowCommitDirty      = false
	_AutoGenerateNextPatch = false
	_AllowDowngrades       = false

	_GenerateChangelog   = true
	_ChangelogFileName   = "CHANGELOG.md"
	_ChangelogTitle      = "Changelog"
	_ChangelogShowAuthor = false
	_ChangelogShowBody   = true

	DefaultConfigFile = "version.yaml"
)
