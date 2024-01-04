package config

// Viper keys.
const (
	AppName = "appName"
	Version = "version"
	WorkDir = "WORK_DIR"

	RemoteURL = "repoURL" // Remote repository URL.

	ConfigFile = "configFile" // Configuration file name. Default: config.yaml.

	// If --force flag is set, then allowCommitDirty, autoGenerateNextPatch and allowDowngrades are set to true.
	AllowCommitDirty      = "allowCommitDirty"      // Allow commit dirty repository. Default: false.
	AutoGenerateNextPatch = "autoGenerateNextPatch" // Auto generate next patch version, if version exists. Default: false.
	AllowDowngrades       = "allowDowngrades"       // Allow version downgrades. Default: false.

	GenerateChangelog   = "changelog.generate"   // Generate changelog. Default: true.
	ChangelogFileName   = "changelog.fileName"   // Changelog file name. Default: CHANGELOG.md.
	ChangelogTitle      = "changelog.title"      // Changelog title. Default: Changelog.
	ChangelogIssueURL   = "changelog.issueURL"   // Issue href template (with last slash). Default: empty.
	ChangelogShowAuthor = "changelog.showAuthor" // Show author in changelog. Default: false.
	ChangelogShowBody   = "changelog.showBody"   // Show body in changelog comment. Default: true.

	ChangelogCommitTypes = "changelog.commitTypes" // Commit types for changelog. Map[string]string. Key - type key, value- type name.
	ChangelogCommitOrder = "changelog.commitOrder" // Commit types order for changelog. []string. Default: empty.

	Backup = "backupChanged" // Backup changed files. Default: false.
	Silent = "silent"        // Silent mode from flags.
	DryRun = "dryRun"        // Dry run mode from flags.
)

const (
	_AppName               = "Version"
	_Version               = "0.0.0"
	_WorkDir               = ""
	_AllowCommitDirty      = false
	_AutoGenerateNextPatch = false
	_AllowDowngrades       = false

	_GenerateChangelog   = true
	_ChangelogFileName   = "CHANGELOG.md"
	_ChangelogTitle      = "Changelog"
	_ChangelogShowAuthor = false
	_ChangelogShowBody   = true

	_ConfigFile = "config.yaml"
)
