![MIT license](https://img.shields.io/badge/license-MIT-green.svg?style=plastic "MIT")
![Version v1.0.5](https://img.shields.io/badge/version-v1.0.5-blue.svg?style=plastic "Version v1.0.5")

[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg)](https://conventionalcommits.org)

![Go](https://img.shields.io/badge/Go-1.21-yellow.svg?style=plastic "Go")

![image info](./img.jpg)

- [Version](#version)
- [Installation](#installation)
- [Usage](#usage)
    - [Common help](#common-help)
    - [Config file](#config-file)
        - [root section](#config-file-root)
            - [version](#config-file-root-version)
            - [backupChanged](#config-file-root-backupChanged)
            - [before and after](#config-file-root-before)
        - [git](#config-file-git)
        - [changelog](#config-file-changelog)
        - [bump files](#config-file-bump)
    - [Changelog format](#changelog-format)
    - [Generate command](#generate-command)
    - [Next command](#next-command)
    - [Remove command](#remove-command)

# <a id='version'>Version</a>

Version is a console utility written in Go, designed for version management in projects following
the [Semantic Versioning 2.0.0](https://semver.org) principles. This tool automates the versioning process in line
with [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0) standards, and assists in generating a
changelog file. It supports version updates in `composer.json`, `package.json`, and other configuration files related to
versioning.

You can try version utility in [version-test](https://github.com/klimby/version-test) repository.

## <a id='installation'>Installation</a>

```bash
wget https://raw.githubusercontent.com/klimby/version/master/bin/version -O version 

sudo chmod +x ./version
```

It is recommended to place the executable in the directory containing your project's git repository.
However, if you need to specify a different project location, you can use the **--dir** flag.
This flag allows setting either an absolute or relative path to the project.

The version utility provides a range of commands for version management and changelog file generation.
When adding the executable file to your project, it is recommended to include it in the `.gitignore` file
to avoid increasing the repository size.

## <a id='usage'>Usage</a>

### <a id='common-help'>Common help</a>

Basic usage:

```bash
$ ./version [command] [flags]
```

```bash
$ ./version --help
CLI tool for versioning, generate changelog, bump version.

Usage:
  version [command]

Available Commands:
  generate    Generate files
  help        Help about any command
  next        Generate next version
  remove      Remove files

Flags:
  -b, --backup          backup changed files
  -c, --config string   config file path (default "version.yaml")
      --dir string      working directory, default - current
  -d, --dry             dry run
  -f, --force           force mode
  -h, --help            help for version
  -s, --silent          silent run
```

Available commands:

* **generate** - Generate config and changelog files.
* **help** - Help about any command.
* **next** - Generate next version.
* **remove** - Remove files.

Global flags (available for all commands):

* **-b**, **--backup** - Backup changed files.

  All original files will be saved with the .bak extension.

  I recommend add `*.bak` to `.gitignore` file.

  You can remove generated backup files with `./version remove --backup` command.

  You can set `backupChanged: true` in config file for enable backup files for all commands.

* **-c**, **--config** (string) - Config file path (default `version.yaml` in working directory).

  You can use `--config` flag for specify config file path. Path can be relative from working directory.

  For example: `./version --config=dir/version.yaml`.

  You can generate config file with `./version generate --config-file` command.

  See details in [Config file](#config-file) section.

* **--dir** (string) - Working directory, default - current.

  You can use `--dir` flag for specify working directory. Path can be absolute or relative from current directory.

  For example: `./version --dir=dir`, `./version --dir=/home/dir`.

* **-d**, **--dry** - Dry run. Files will not be changed, commit will not be created.

* **-f**, **--force** - Force mode. When enabled:
    - Files will be changed, even if the repository is not clean (see [Config file](#config-file) `git.commitDirty`
      parameter).
    - If the version already exists, next patch will be generated (see [Config file](#config-file) `git.autoNextPatch`
      parameter).
    - Allow version downgrades with `--ver` flag (see [Config file](#config-file) `git.allowDowngrades`
      parameter).

* **-h**, **--help** - Help for command.
* **-s**, **--silent** - Silent run. No output. If you use this flag, then you will not see any output from the
  utility. This is useful if you want to use the utility in scripts. If app finished with error, will be returned exit
  code 1, else 0.

### <a id='config-file'>Config file</a>

Although the version utility can operate without a configuration file using default values,
customization is often necessary for specific project requirements.

To create a configuration file, use the command `version generate --config-file`.

This will generate a `version.yaml` file or a custom configuration file if the **--config** flag is specified.

The configuration file includes:

* Backup settings for changed files.
* Commands to be executed before and after commit.
* Git settings, including parameters for automatic patch creation and allowing version downgrades.
* Settings for generating the changelog file.
* Rules for changing the version in specific files, with the option to use regular expressions for search and replace.

Example:

```yaml
# Version configuration file.
# Generated by Version v0.0.1.

# Application version.
version: 1.0.0

# Backup changed files.
# All original files will be saved with the .bak extension.
backupChanged: false

# Run commands before commit.
# All commands will be executed from main directory (where version is located).
# Parameters:
#   - cmd: command for run in format: [ "echo", "after commit" ]
#   - versionFlag: flag to send bumped version in format 1.2.3 to command. Optional.
#   - breakOnError: flag that indicates that the command is stopped if an error occurs. Optional.
#   - runInDry: flag that indicates that the command is run in dry mode. Optional.
# Examples:
# before:
#   - cmd: [ "echo", "before commit" ]
#     versionFlag: "--version"
#     breakOnError: true
#     runInDry: true
# In this example, will be run command: echo before commit --version=1.2.3
before:
  - cmd: [ "make", "build" ]
    versionFlag: "VERSION"
    breakOnError: true
    runInDry: false
  - cmd: [ "cp", "./bin/version", "./" ]
    breakOnError: true
    runInDry: false

# Run commands after commit.
# All commands will be executed from main directory (where version is located).
# Parameters:
# 	- cmd: command for run in format: [ "echo", "after commit" ]
# 	- versionFlag: flag to send bumped version in format 1.2.3 to command. Optional.
# 	- breakOnError: flag that indicates that the command is stopped if an error occurs. Optional.
# 	- runInDry: flag that indicates that the command is run in dry mode. Optional.
# Examples:
# after:
# 	- cmd: [ "echo", "after commit" ]
# 	  versionFlag: "--version"
# 	  breakOnError: true
# 	  runInDry: true
# In this example, will be run command: echo after commit --version=1.2.3
after:
  - cmd: [ "echo", "after commit" ]

# Git settings.
git:
  # Allow commit not clean repository.
  commitDirty: false
  # Auto generate next patch version, if version exists.
  autoNextPatch: false
  # Allow version downgrades with --version flag.
  allowDowngrades: false
  # Remote repository URL.
  remoteUrl: https://github.com/klimby/version

# Changelog settings.
changelog:
  # Generate changelog.
  generate: true
  # Changelog file name.
  file: CHANGELOG.md
  # Changelog title.
  title: "Changelog"
  # Issue url template.
  # Examples:
  #  - IssueURL: https://company.atlassian.net/jira/software/projects/PROJECT/issues/
  #  - IssueURL: https://github.com/company/project/issues/
  # If empty, ang repository is CitHub, then issueHref will be set from remote repository URL.
  issueUrl: https://github.com/klimby/version/issues/
  # Show author in changelog.
  showAuthor: false
  # Show body in changelog comment.
  showBody: true
  # Commit types for changelog.
  # Type - commit type, value - commit type name.
  # If empty, then all commit types will be used, except Breaking Changes.
  commitTypes:
    - type: "feat"
      name: "Features"
    - type: "fix"
      name: "Bug Fixes"
    - type: "perf"
      name: "Performance Improvements"
    - type: "refactor"
      name: "Code Refactoring"
    - type: "style"
      name: "Styles"
    - type: "test"
      name: "Tests"
    - type: "build"
      name: "Builds"
    - type: "docs"
      name: "Documentation"
    - type: "revert"
      name: "Reverts"
    - type: "ci"
      name: "Continuous Integration"
    - type: "chore"
      name: "Other changes"

# Bump files.
# Change version in files. Version will be changed with format: <digital>.<digital>.<digital>
# Every entry has format:
# - file: file patch
#   start: number (optional)
#   end: number (optional)
#   regexp: regular expression for string search (optional, array)
# 
# If file is composer.json or package.json, then regexp and start/end are ignored.
#
# Examples:
# bump:
#   - file: README.md
#     regexp: 
#       - ^Version:.+$
# 
# All strings from file, that match regexp will be replaced with new version.
#
# bump:
#   - file: dir/file.txt
#     start: 0
#     end: 100
#
# All strings from file, from 0 to 100 will be replaced with new version.
#
# bump:
#   - file: README.md
#     regexp: 
#       - ^Version:.+$
#     start: 0
#     end: 100
#
# All strings from file, from 0 to 100, that match regexp will be replaced with new version.
#
bump:
  - file: package.json
  - file: README.md
    start: 0
    end: 5
    regexp:
      - ^!\[Version:.*$
  - file: foo/bar.txt
    start: 0
    end: 5

```

#### <a id='config-file-root'>root section</a>

##### <a id='config-file-root-version'>version</a>

Application version. Don't change this value. It needs for check config file version for future updates.

##### <a id='config-file-root-backupChanged'>backupChanged</a>

Backup changed files. All original files will be saved with the .bak extension.

You can remove generated backup files with `./version remove --backup` command.

You can set flag * **-b**, **--backup** for enable backup files for specific command
(see [Common help](#common-help) section).

##### <a id='config-file-root-before'>before and after</a>

Run commands before and after commit. All commands will be executed from main directory (where version is located).

Parameters:

* **cmd** - command for run in format: `[ "echo", "after commit" ]`.
* **versionFlag** - flag to send bumped version in format 1.2.3 to command. Optional.

  For example, if you set `versionFlag: "--version"`, then command will be run with flag `--version=1.2.3`.

* **breakOnError** - flag that indicates that the command is stopped if an error occurs. Optional.
* **runInDry** - flag that indicates that the command is run in dry mode. Optional.

Examples:

```yaml
before:
  - cmd: [ "echo", "before commit" ]
    versionFlag: "--version"
    breakOnError: true
    runInDry: true
```

In this example, will be run command: `echo before commit --version=1.2.3`.

#### <a id='config-file-git'>git</a>

Git and commit settings.

* **commitDirty** - allow commit not clean repository.
* **autoNextPatch** - auto generate next patch version, if version exists.
* **allowDowngrades** - allow version downgrades with `--ver` flag.
* **remoteUrl** - remote repository URL. For GitHub repository it sets from remote repository URL as default.

If you run command next with **--force** flag, then:

* If **commitDirty** is true, then commit will be created, even if the repository is not clean.
* If **autoNextPatch** is true, then next patch version will be generated, if version exists.
* If **allowDowngrades** is true, then version downgrade will be allowed.

#### <a id='config-file-changelog'>changelog</a>

Changelog settings. See also [Changelog format](#changelog-format) section.

* **generate** - generate changelog.
* **file** - changelog file name. You can use relative path from working directory.
* **title** - changelog title (first line).
* **issueUrl** - issue url template.

  If repository is CitHub, then issueHref will be set from remote repository URL as default.

  For example, for JIRA you can set `issueUrl: https://company.atlassian.net/jira/software/projects/PROJECT/issues/`.

* **showAuthor** - show commit author in changelog.
* **showBody** - show commit body in changelog comment.
* **commitTypes** - commit types for changelog.

  Type - commit type, value - commit type name for markdown.

  If empty, then all commit types will be hidden, except Breaking Changes.

  For example:

  ```yaml
  commitTypes:
    - type: "feat"
      name: "Features"
    - type: "fix"
      name: "Bug Fixes"
    - type: "perf"
      name: "Performance Improvements"
    - type: "refactor"
      name: "Code Refactoring"
    - type: "style"
      name: "Styles"
    - type: "test"
      name: "Tests"
    - type: "build"
      name: "Builds"
    - type: "docs"
      name: "Documentation"
    - type: "revert"
      name: "Reverts"
    - type: "ci"
      name: "Continuous Integration"
    - type: "chore"
      name: "Other changes"
  ```

#### <a id='config-file-bump'>bump files</a>

Change version in files. Version will be changed with format: `<digital>.<digital>.<digital>`.

Every entry has format:

* **file** - file patch. Relative path from working directory. Required.
* **start** - start line number. Optional.
* **end** - end line number. Optional.
* **regexp** - regular expressions for string search. Optional, array.

  Must contains valid Go regexps.

If file is `composer.json` or `package.json`, then regexp and start/end are ignored.

Common behavior:

1. Find target file.
2. If file is `composer.json` or `package.json`, then version will be changed in `version` field.
3. If file is not `composer.json` or `package.json`, then:
    1. If start/end is set, then analyze only lines from start to end. Else analyze all lines.
    2. If regexp are set, then analyze only lines from 3.1 that match one of regexp. Else analyze all lines from 3.1.
    3. If line contains version in format `<digital>.<digital>.<digital>`, then version will be changed in this line.

If regexp is invalid, then will be returned error on start of app.

If file not found, then will be returned error on start of app.

### <a id='changelog-format'>Changelog format</a>

Changelog will be generated in **markdown** format.

Changelog generate with parameters from [changelog](#config-file-changelog) section in [config file](#config-file).

Commits must be in the format [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0).

Changelog has title, defined in [config file](#config-file-changelog) `title` parameter.

Every version has format:

```markdown
## [1.2.3](compare-url(1)) (Date in format 2021-01-01)

### Breaking changes(2)

* commits(3)

### Features(4)

* commits(3)

### Bug fixes(4)

...

```

1. Compare URL - compare URL between previous and current version. If previous version is empty, then compare URL will
   be empty.
2. Breaking changes - all breaking changes commits.
   See [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0).
3. Commits - all commits for this type. Commit format see below.
4. Commit types - commit types from [config file](#config-file-changelog) `commitTypes` parameter in order from config
   file.

All commits without type will be added to **Other changes** section with type **chore**.

If commit has **BREAKING CHANGE** in body or ! in type, then commit will be added to **Breaking changes** section.

If `commitTypes` parameter is empty, then all commit types will be hidden, except Breaking Changes.

Commit format:

```markdown
* **scope(1):** commit message ([commit hash(2)](commit url(3))) ([author(4)](author url(5)))
    * commit body(6)
    * commit body with issue [123](issue-url)(7)
```

1. Scope - commit scope, if exists.
2. Commit hash - commit short hash.
3. Commit URL - commit URL, if `remoteUrl` parameter is set in [config file](#config-file-git).
4. Author - commit author, if `showAuthor` parameter is true in [config file](#config-file-changelog).
5. Author URL - author URL, if `remoteUrl` parameter is set in [config file](#config-file-git).

   Url calculated from git commit author email.

6. Commit body - commit body, if `showBody` parameter is true in [config file](#config-file-changelog).
7. Issue in body must be start from #. For example, #123. If issue URL is set in [config file](#config-file-changelog),
   then issue will be linked to issue URL.

For example, you can use `CHANGELOG.md` file in this project.

### <a id='generate-command'>Generate command</a>

Generate full changelog and config file:

```bash
$ version generate --help
CLI tool for versioning, generate changelog, bump version.

Usage:
  version generate [flags]

Flags:
      --changelog     generate changelog file
      --config-file   generate config file
  -h, --help          help for generate

Global Flags:
  -b, --backup          backup changed files
  -c, --config string   config file path (default "version.yaml")
      --dir string      working directory, default - current
  -d, --dry             dry run
  -f, --force           force mode
  -s, --silent          silent run
  -v, --verbose         verbose output
```

* **--changelog** - generate changelog file. If file exists, then will be rewritten.
* **--config-file** - generate config file. If file exists, then will be rewritten.

### <a id='next-command'>Next command</a>

Command for creating next version, add content to changelog, bump files and commit changes:

```bash
$ version next --help
Generate next version.

Usage:
  version next [flags]

Flags:
  -h, --help         help for next
      --major        next major version
      --minor        next minor version
      --patch        next patch version
      --prepare      run only bump files and commands before
      --ver string   next build version in format 1.2.3

Global Flags:
  -b, --backup          backup changed files
  -c, --config string   config file path (default "version.yaml")
      --dir string      working directory, default - current
  -d, --dry             dry run
  -f, --force           force mode
  -s, --silent          silent run
  -v, --verbose         verbose output
```

* **--major** - next major version.
* **--minor** - next minor version.
* **--patch** - next patch version.
* **--ver** - next build version in format 1.2.3. For example: `--ver=1.2.3`.

### <a id='remove-command'>Remove command</a>

Command for remove backup files:

```bash
$ version remove --help
Remove backup files, if exists

Usage:
  version remove [flags]

Flags:
      --backup   remove backup files
  -h, --help     help for remove

Global Flags:
  -c, --config string   config file path (default "version.yaml")
      --dir string      working directory, default - current
  -d, --dry             dry run
  -f, --force           force mode
  -s, --silent          silent run
```

* **--backup** - remove backup files.
