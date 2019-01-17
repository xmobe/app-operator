{{- define "componentName" }}track{{ end }}
{{- define "componentType" }}operator{{ end }}
apiVersion: cluster.odoo.io/v1beta1
kind: OdooTrack
metadata:
  name: {{ .Extra.Name }}.{{ block "componentType" . }}{{ end }}.{{ block "componentName" . }}{{ end }}
  namespace: {{ .Instance.Namespace }}
spec:
  track:  {{ .Extra.Track | quote }}
