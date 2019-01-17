{{- define "componentName" }}config{{ end }}
{{- define "componentType" }}app{{ end }}
apiVersion: v1
kind: ConfigMap
{{- template "metadata" . -}}
data:
  app-config: |
{{ .Extra.ConfigFile | indent 4}}
