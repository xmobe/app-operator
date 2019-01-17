{{- define "metadata"}}
metadata: &metadata
  name: {{ .Instance.Name }}.{{ block "componentType" . }}{{ end }}.{{ block "componentName" . }}{{ end }}
  namespace: {{ .Instance.Namespace }}
  labels: &metadatalabels
    cluster.odoo.io/part-of-cluster: {{ .Instance.Spec.Cluster | quote }}
    instance.odoo.io/hostname: {{ .Instance.Spec.Hostname | quote }}
    app.kubernetes.io/name: {{ block "componentName" . }}{{ end }}
    app.kubernetes.io/instance: {{ .Instance.Name }}-{{ block "componentName" . }}{{ end }}
    app.kubernetes.io/component: {{ block "componentType" . }}{{ end }}
    app.kubernetes.io/managed-by: odoo-operator
    app.kubernetes.io/part-of: {{ .Instance.Name | quote }}
    app.kubernetes.io/version: {{ .Instance.Spec.Version | quote }}
    app.kubernetes.io/track: {{ .Extra.Track | quote }}
{{ end -}}