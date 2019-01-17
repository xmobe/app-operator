{{- define "metadata" }}
metadata:
  name: {{ .Instance.Name }}.{{ block "componentType" . }}{{ end }}.{{ block "componentName" . }}{{ end }}
  namespace: {{ .Instance.Namespace }}
  labels:
    app.kubernetes.io/name: {{ block "componentName" . }}{{ end }}
    app.kubernetes.io/instance: {{ .Instance.Name }}-{{ block "componentName" . }}{{ end }}
    app.kubernetes.io/component: {{ block "componentType" . }}{{ end }}
    app.kubernetes.io/managed-by: odoo-operator
    app.kubernetes.io/part-of: {{ .Instance.Name | quote }}
    app.kubernetes.io/version: NA
    app.kubernetes.io/track: NA
{{ end -}}