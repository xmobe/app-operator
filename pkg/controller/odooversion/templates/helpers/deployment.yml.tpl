{{- define "deployment" }}
apiVersion: apps/v1
kind: Deployment
{{- template "metadata" . -}}
spec:
  replicas: {{ block "deploymentReplicas" . }}1{{ end }}
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 1
      maxSurge: 1
  selector:
    matchLabels:
      app.kubernetes.io/instance: {{ .Instance.Name }}-{{ block "componentName" . }}{{ end }}
  template:
    metadata:
      labels:
        app.kubernetes.io/name: {{ block "componentName" . }}{{ end }}
        app.kubernetes.io/instance: {{ .Instance.Name }}-{{ block "componentName" . }}{{ end }}
        app.kubernetes.io/version: {{ .Instance.Spec.Version }}
        app.kubernetes.io/component: {{ block "componentType" . }}{{ end }}
        app.kubernetes.io/part-of: {{ .Instance.Name }}
        app.kubernetes.io/managed-by: odoo-operator
    spec:
      imagePullSecrets:
      - name: pull-secret
      terminationMessagePolicy: FallbackToLogsOnError
      securityContext:
        fsGroup: 9001
        runAsUser: 9001
        runAsGroup: 9001
        runAsNonRoot: true
        supplementalGroups: [2000]
      containers:
      - name: default
        image: {{ .Extra.Image }}:base-{{ .Instance.Spec.Track }}-{{ .Instance.Spec.Version }}
        imagePullPolicy: Always
        args:
        {{ block "deploymentArgs" . }}{{ end }}
        ports: {{ block "deploymentPorts" . }}[]{{ end }}
        {{ block "deploymentHealchecks" . }}{{ end }}
        resources:
          requests:
            memory: 512M
            cpu: 200m
          limits:
            memory: 1G
            cpu: 500m
        env:
         - name: PGHOST
           value: {{ .Extra.Database.Host }}
         - name: PGUSER
           value: {{ .Extra.Database.User }}
         - name: PGPORT
           value: {{ .Extra.Database.Port }}
         - name: PGPASSWORD
           value: {{ .Extra.Database.Password }}
         - name: ODOO_RC
           value: /run/configs/odoo/
         - name: ODOO_PASSFILE
           value: /run/secrets/odoo/adminpwd
        volumeMounts:
        - name: data-volume
          mountPath: /mnt/odoo/data/
        - name: config-volume
          mountPath: /run/configs/odoo/
          readonly: true
        - name: app-secrets
          mountPath: /run/secrets/odoo/
          readonly: true
      volumes:
        - name: data-volume
          volumeSource:
            name: {{ .Extra.ClusterName }}-data
        - name: backup-volume
          volumeSource:
            name: {{ .Extra.ClusterName }}-backup
        - name: config-volume
          configMap:
            name: {{ .Instance.Name }}-config
            defaultMode: 272
        - name: app-secrets
          secret:
            secretName: {{ .Instance.Name }}.app-secrets
            defaultMode: 256
{{ end -}}
