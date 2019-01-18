{{- define "componentName" }}initializer{{ end }}
{{- define "componentType" }}db{{ end }}
{{- define "jobArgs" -}}
        - dodoo-initializer
        - --config
        - /run/configs/odoo/
        - --new-database
        - {{ .Instance.Spec.Hostname }}
	{{- if .Instance.Spec.InitModules }}
        - --modules
        - {{ .Instance.Spec.InitModules | join "," }}
	{{- end -}}
{{- end -}}
apiVersion: batch/v1
kind: Job
{{- template "metadata" . -}}
{{- template "jobspec" . -}}
