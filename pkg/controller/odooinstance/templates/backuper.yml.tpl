{{- define "componentName" }}backuper{{ end }}
{{- define "componentType" }}db{{ end }}
{{- define "jobArgs" -}}
        - dodoo
        - snapshot
        - --config
        - /run/configs/odoo/
        - {{ .Instance.Spec.Hostname }}
        - destination-folder
{{- end -}}
apiVersion: batch/v1beta1
kind: CronJob
{{- template "metadata" . -}}
{{- template "jobspec" . -}}
spec:
  schedule: "* * * 1 *"
  concurrencyPolicy: Forbid
  jobTemplate:
    matadata: *metadata
    spec: *jobspec
