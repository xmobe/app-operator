{{- define "service" }}
kind: Service
apiVersion: v1
{{- template "metadata" . -}}
spec:
  selector: *metadatalabels
  type: ClusterIP
  ports:
  {{ block "servicePorts" . }}{{ end }}
{{ end -}}