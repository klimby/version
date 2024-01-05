package config

const _configYamlTemplate = `# Version configuration file.
# Generated by Version {{ .Version }}.

# Application settings.
version: {{ .Version }}

# Backup changed files.
# All original files will be saved with the .bak extension.
backupChanged: {{ .Backup }}

# Git settings.
git:
	# Allow commit not clean repository.
	commitDirty: {{ .GitOptions.AllowCommitDirty }}
	# Auto generate next patch version, if version exists.
	autoNextPatch: {{ .GitOptions.AutoGenerateNextPatch }}
	# Allow version downgrades with --version flag.
	allowDowngrades: {{ .GitOptions.AllowDowngrades }}
	# Remote repository URL.
	remoteUrl: {{ .GitOptions.RemoteURL }}

# Changelog settings.
changelog:
	# Generate changelog.
	generate: {{ .ChangelogOptions.Generate }}
	# Changelog file name.
	file: {{ .ChangelogOptions.FileName.String }}
	# Changelog title.
	title: "{{ .ChangelogOptions.Title }}"
	# Issue url template.
	# Examples:
	# 	- IssueURL: https://company.atlassian.net/jira/software/projects/PROJECT/issues/
	# 	- IssueURL: https://github.com/company/project/issues/
	# If empty, ang repository is CitHub, then issueHref will be set from remote repository URL.
	issueUrl: {{ .ChangelogOptions.IssueURL }}
	# Show author in changelog.
	showAuthor: {{ .ChangelogOptions.ShowAuthor }}
	# Show body in changelog comment.
	showBody: {{ .ChangelogOptions.ShowBody }}
	# Commit types for changelog.
	# Type - commit type, value - commit type name.
	# If empty, then all commit types will be used, except Breaking Changes.
	commitTypes:
	{{- range .ChangelogOptions.CommitTypes }}
		- type: "{{ .Type }}"
		  name: "{{ .Name }}"
	{{- end}}

# Bump files.
# Change version in files. Version will be changed with format: <digital>.<digital>.<digital>
# Every entry has format:
# - file: file patch
#   start: number (optional)
#   end: number (optional)
#   regexp: regular expression for string search (optional)
# 
# If file is composer.json or package.json, then regexp and start/end are ignored.
#
# Examples:
# - file: README.md
#   regexp: 
#		- ^Version:.+$
# 
# All strings from file, that match regexp will be replaced with new version.
#
# - file: dir/file.txt
#   start: 0
#   end: 100
#
# All strings from file, from 0 to 100 will be replaced with new version.
#
# - file: README.md
#   regexp: 
#		- ^Version:.+$
#   start: 0
#   end: 100
#
# All strings from file, from 0 to 100, that match regexp will be replaced with new version.
#
bump: 
{{- range $value := .Bump }}
	- file: {{ $value.File.String }}
	{{- if $value.HasPositions }}
	  start: {{ $value.Start }}
	  end: {{ $value.End }}
	{{- end}}
	{{- if $value.HasRegExp }}
	  regexp: 
	  {{- range $value.RegExp }}
	    - {{ . }}
      {{- end}}
	{{- end}}
{{- end}}	
`
