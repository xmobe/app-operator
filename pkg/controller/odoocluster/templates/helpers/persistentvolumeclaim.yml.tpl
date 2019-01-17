{{- define "pvc" }}
apiVersion: v1
kind: PersistentVolumeClaim
{{- template "metadata" . -}}
spec:
  accessModes:
    - ReadWriteOnce
  volumeMode: Filesystem
  resources:
    requests:
      storage: 3Gi
{{ end -}}
