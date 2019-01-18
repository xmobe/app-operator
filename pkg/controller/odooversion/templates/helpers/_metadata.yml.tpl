{{- define "metadata" }}
metadata:
  {{- if .Extra.Service }}
  name: v{{ .Instance.Spec.Version | replace "." "-" }}-{{ block "componentType" . }}{{ end }}-{{ block "componentName" . }}{{ end }}
  {{- else }}
  name: v{{ .Instance.Spec.Version | replace "." "-" }}.{{ block "componentType" . }}{{ end }}.{{ block "componentName" . }}{{ end }}
  {{- end }}
  namespace: {{ .Instance.Namespace }}
  labels: &metadatalabels
    cluster.odoo.io/part-of-cluster: {{ .Instance.Spec.Cluster | quote }}
    app.kubernetes.io/name: {{ block "componentName" . }}{{ end }}
    app.kubernetes.io/instance: {{ .Instance.Name }}-{{ block "componentName" . }}{{ end }}
    app.kubernetes.io/component: {{ block "componentType" . }}{{ end }}
    app.kubernetes.io/managed-by: odoo-operator
    app.kubernetes.io/part-of: {{ .Instance.Name | quote }}
    app.kubernetes.io/version: {{ .Instance.Spec.Version | quote }}
    app.kubernetes.io/track: {{ .Instance.Spec.Track | quote }}
{{ end -}}