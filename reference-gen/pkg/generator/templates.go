package generator

import (
	"fmt"
	"path/filepath"
	"text/template"
)

var defaultTemplates = []string{
	packageTemplate,
	typeTemplate,
	memberTemplate,
	membersTemplate,
	memberWithEmbedTemplate,
}

const packageTemplate = `
{{- define "package" -}}
    {{- range (visibleTypes (sortedTypes .types)) -}}
        {{ template "type" .  }}
    {{- end -}}
{{- end -}}
`

const typeTemplate = `
{{ define "type" }}
### {{ .Name.Name }}
{{- if or (eq .Kind "Alias") (aliasDisplayName .) }}
{{ if linkForType .Underlying }}
#### ([{{ aliasDisplayName . }}]({{ linkForType .Underlying}}) alias)
{{- else -}}
#### ({{ backtick (aliasDisplayName .)}} alias)
{{- end -}}
{{ end }}
{{ with (typeReferences .) }}
(**Appears on:** {{ $prev := "" -}}
    {{- range . -}}
        {{- if $prev -}}, {{ end -}}
        {{- $prev = . -}}
        [{{ typeDisplayName . }}]({{ linkForType . }})
    {{- end -}}
  )
{{ end }}
{{ renderCommentsLF .CommentLines }}
{{ if visibleMembers .Members }}
| Field | Type | Description |
| ----- | ---- | ----------- |
{{- template "members_with_embed" . }}
{{ end -}}
{{ end }}
`

const memberTemplate = `
{{ define "member" }}
  {{- if not (hideMember .) }}
| {{ backtick (fieldName .) }} | _{{- if linkForType .Type -}}
    [{{ typeDisplayName .Type }}]({{ linkForType .Type}})
  {{- else -}}
    {{ typeDisplayName .Type }}
  {{- end -}}_ | {{ if fieldEmbedded . -}}
    (Members of {{ backtick (fieldName .) }} are embedded into this type.)
  {{ end -}}
  {{- if isOptionalMember . }} _(Optional)_ {{ end -}}
  {{- renderCommentsBR .CommentLines }} |
  {{- end -}}
{{- end }}
`

const membersTemplate = `
{{ define "members" }}
{{- range .Members -}}
{{- template "member" . -}}
{{- end -}}
{{ end }}
`

const memberWithEmbedTemplate = `
{{ define "members_with_embed" }}
{{- range .Members -}}
{{- if fieldEmbedded . -}}
{{- template "members" (dereference .Type) -}}
{{- else -}}
{{- template "member" . -}}
{{- end -}}
{{- end -}}
{{ end }}
`

// loadTemplatesInto loads templates from the directory given, or the default
// templates, into the template object.
func loadTemplatesInto(t *template.Template, templateDir string) (*template.Template, error) {
	// No template directory given, use default templates
	if templateDir == "" {
		return loadDefaultTemplatesInto(t)
	}

	return t.ParseGlob(filepath.Join(templateDir, "*.tpl"))
}

// loadDefaultTemplatesInto loads the defaultTemplates into the template object.
func loadDefaultTemplatesInto(t *template.Template) (*template.Template, error) {
	var err error
	for _, tmpl := range defaultTemplates {
		t, err = t.Parse(tmpl)
		if err != nil {
			return nil, fmt.Errorf("error loading default template: %v", err)
		}
	}
	return t, nil
}
