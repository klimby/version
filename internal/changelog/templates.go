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
{{- if .CommitsFromLast }}

### {{ .Name }}
{{ range .CommitsFromLast}}
* {{ commitName . }}
{{- range .Body}}
    * {{addIssueURL .}}
{{- end}}
{{- end}}
{{- end -}}
{{- end}}
`
