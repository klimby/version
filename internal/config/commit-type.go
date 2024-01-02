package config

const (
	CommitFeat     = "feat"     // Features
	CommitFix      = "fix"      // Bug Fixes
	CommitPerf     = "perf"     // Performance Improvements
	CommitRefactor = "refactor" // Code Refactoring
	CommitStyle    = "style"    // Styles
	CommitTest     = "test"     // Tests
	CommitBuild    = "build"    // Builds
	CommitChore    = "chore"    // Other changes
	CommitDocs     = "docs"     // Documentation
	CommitRevert   = "revert"   // Reverts
	CommitCI       = "ci"       // Continuous Integration
)

type CommitName struct {
	Type string `yaml:"type"`
	Name string `yaml:"name"`
}

var _defaultCommitNames = []CommitName{
	{Type: CommitFeat, Name: "Features"},
	{Type: CommitFix, Name: "Bug Fixes"},
	{Type: CommitPerf, Name: "Performance Improvements"},
	{Type: CommitRefactor, Name: "Code Refactoring"},
	{Type: CommitStyle, Name: "Styles"},
	{Type: CommitTest, Name: "Tests"},
	{Type: CommitBuild, Name: "Builds"},
	{Type: CommitDocs, Name: "Documentation"},
	{Type: CommitRevert, Name: "Reverts"},
	{Type: CommitCI, Name: "Continuous Integration"},
	{Type: CommitChore, Name: "Other changes"},
}

func toViperCommitNames(names []CommitName) (types map[string]string, order []string) {
	types = make(map[string]string, len(names))
	order = make([]string, len(names))

	for i, name := range names {
		types[name.Type] = name.Name
		order[i] = name.Type
	}

	return types, order
}

func fromViperCommitNames(types map[string]string, order []string) []CommitName {
	names := make([]CommitName, 0, len(types))

	for _, tp := range order {
		names = append(names, CommitName{
			Type: tp,
			Name: types[tp],
		})
	}

	return names
}
