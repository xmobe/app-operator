{{- define "componentName" }}syncmigrator{{ end }}
{{- define "componentType" }}db{{ end }}
{{- define "jobArgs" -}}
        - dodoo-migrator
        - --config
        - /run/configs/odoo/
        - --database
        - {{ .Instance.Spec.Hostname }}
        - --file
        - /opt/odoo/.migration.yml
{{- end -}}
apiVersion: batch/v1
kind: Job
{{- template "metadata" . -}}
{{- template "jobspec" . -}}
