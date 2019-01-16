{{- define "metadata" }}
metadata:
  name: {{ block "componentType" . }}{{ end }}-{{ .Instance.Spec.Version | replace "." "-" }}-{{ block "componentName" . }}{{ end }}
  namespace: {{ .Instance.Namespace }}
  labels: &metadatalabels
    cluster.odoo.io/name: {{ .Instance.Spec.Cluster | quote }}
    cluster.odoo.io/track: {{ .Instance.Spec.Track | quote }}
    app.kubernetes.io/name: {{ block "componentName" . }}{{ end }}
    app.kubernetes.io/instance: {{ .Instance.Name }}-{{ block "componentName" . }}{{ end }}
    app.kubernetes.io/component: {{ block "componentType" . }}{{ end }}
    app.kubernetes.io/managed-by: odoo-operator
    app.kubernetes.io/part-of: {{ .Instance.Name | quote }}
    app.kubernetes.io/version: {{ .Instance.Spec.Version | quote }}
{{ end -}}