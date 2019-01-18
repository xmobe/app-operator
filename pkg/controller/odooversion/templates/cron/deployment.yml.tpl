{{ $service := false }}
{{- define "componentName" }}cron{{ end }}
{{- define "componentType" }}app{{ end }}
{{- define "deploymentArgs" -}}
        - --config
        - /run/configs/odoo/
        - --db_maxconn=1
        - --workers=0
        - --max-cron-threads=1
        - --no-xmlrpc
{{- end -}}
{{- template "deployment" . -}}