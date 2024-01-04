package changelog

const _tagMarkdownTpl = `
## {{ versionName . }} ({{.Date}})
{{- if .BreakingChanges}}

### Breaking changes
{{- range .BreakingChanges}}

* {{ commitName . }}
{{- range .Body}}
    * {{addIssueURL .}}
{{- end}}
{{- end}}
{{- end -}}

{{- range .Blocks}}
{{- if .Commits }}

### {{ .Name }}
{{ range .Commits}}
* {{ commitName . }}
{{- range .Body}}
    * {{addIssueURL .}}
{{- end}}
{{- end}}
{{- end -}}
{{- end}}
`
