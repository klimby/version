package config

// Commit types.
const (
	_CommitFeat     = "feat"     // Features
	_CommitFix      = "fix"      // Bug Fixes
	_CommitPerf     = "perf"     // Performance Improvements
	_CommitRefactor = "refactor" // Code Refactoring
	_CommitStyle    = "style"    // Styles
	_CommitTest     = "test"     // Tests
	_CommitBuild    = "build"    // Builds
	CommitChore     = "chore"    // Other changes
	_CommitDocs     = "docs"     // Documentation
	_CommitRevert   = "revert"   // Reverts
	_CommitCI       = "ci"       // Continuous Integration
)

// commitName is a commit type name.
type commitName struct {
	Type string `yaml:"type"`
	Name string `yaml:"name"`
}

var _defaultCommitNames = []commitName{
	{Type: _CommitFeat, Name: "Features"},
	{Type: _CommitFix, Name: "Bug Fixes"},
	{Type: _CommitPerf, Name: "Performance Improvements"},
	{Type: _CommitRefactor, Name: "Code Refactoring"},
	{Type: _CommitStyle, Name: "Styles"},
	{Type: _CommitTest, Name: "Tests"},
	{Type: _CommitBuild, Name: "Builds"},
	{Type: _CommitDocs, Name: "Documentation"},
	{Type: _CommitRevert, Name: "Reverts"},
	{Type: _CommitCI, Name: "Continuous Integration"},
	{Type: CommitChore, Name: "Other changes"},
}

// toViperCommitNames converts commit names to viper types.
func toViperCommitNames(names []commitName) (types map[string]string, order []string) {
	types = make(map[string]string, len(names))
	order = make([]string, len(names))

	for i, name := range names {
		types[name.Type] = name.Name
		order[i] = name.Type
	}

	return types, order
}

// fromViperCommitNames converts viper types to commit names.
func fromViperCommitNames(types map[string]string, order []string) []commitName {
	names := make([]commitName, 0, len(types))

	for _, tp := range order {
		names = append(names, commitName{
			Type: tp,
			Name: types[tp],
		})
	}

	return names
}
