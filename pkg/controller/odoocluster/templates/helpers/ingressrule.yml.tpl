{{- define "ingressrule" }}
  - host: {{ .Spec.Hostname }}
    http:
      paths:
      - path: /longpolling
        backend:
          serviceName: v{{ .Spec.Version | replace "." "-" }}-app-longpolling
          servicePort: 8072
      - path: /
        backend:
          serviceName: v{{ .Spec.Version | replace "." "-" }}-app-web
          servicePort: 8069
{{ end -}}