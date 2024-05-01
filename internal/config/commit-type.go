// Package config provides configuration.
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

// CommitName is a commit type name.
type CommitName struct {
	Type string `yaml:"type"`
	Name string `yaml:"name"`
}

var _defaultCommitNames = []CommitName{
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
