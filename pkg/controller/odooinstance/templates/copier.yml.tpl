{{- define "componentName" }}copier{{ end }}
{{- define "componentType" }}db{{ end }}
{{- define "jobArgs" -}}
        - dodoo
        - copy
        - --config
        - /run/configs/odoo/
        - --force-disconnect
        - {{ .Instance.Spec.ParentHostname }}
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
