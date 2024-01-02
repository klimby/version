package changelog

type TemplateType uint

const (
	MarkdownTpl TemplateType = iota + 1
	ConsoleTpl
)

const _tagMarkdownTpl = `
## {{ versionName . }} ({{.Date}})

{{- if .BreakingChanges}}

### Breaking changes
{{ range .BreakingChanges}}
* {{ commitName . }}
{{- range .Body}}
    * {{addIssueURL .}}
{{- end}}
{{- end}}
{{end -}}
{{- range .Blocks}}
{{- if .Commits }}
### {{ .Name }}
{{- range .Commits}}

* {{ commitName . }}
{{- range .Body}}
    * {{addIssueURL .}}
{{- end }}
{{end -}}
{{end -}}
{{end -}}
`

const _tagConsoleTpl = `
## {{ versionName . }} ({{.Date}})

{{- if .BreakingChanges}}

### Breaking changes
{{ range .BreakingChanges}}
* {{ commitName . }}
{{- range .Body}}
    * {{addIssueURL .}}
{{- end}}
{{- end}}
{{end -}}
`
