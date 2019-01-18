{{- define "jobspec" }}
spec: &jobspec
  completions: 1
  backoffLimit: 1
  activeDeadlineSeconds: 360
  template:
    metadata:
      labels: *metadatalabels
    spec:
      restartPolicy: Never
      imagePullSecrets:
      - name: {{ .Extra.Image.Secret }}
      terminationMessagePolicy: FallbackToLogsOnError
      securityContext:
        fsGroup: 9001
        runAsUser: 9001
        runAsGroup: 9001
        runAsNonRoot: true
        supplementalGroups: [2000]
      containers:
      - name: default
        image: {{ .Extra.Image.Registry }}/{{ .Extra.Image.Repository }}:devops-{{ .Extra.Image.Trackname }}-{{ .Extra.Image.Version }}
        imagePullPolicy: Always
        args:
        {{ block "jobArgs" . }}{{ end }}
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
        - name: app-secret
          mountPath: /run/secrets/odoo/
          readonly: true

      volumes:
        - name: data-volume
          persistentVolumeClaim:
            claimName: {{ .Instance.Spec.Cluster }}.storage.data
        - name: config-volume
          configMap:
            name: v{{ .Instance.Spec.Version | replace "." "-" }}.app.config
        - name: app-secret
          secret:
            secretName: {{ .Instance.Spec.Cluster }}.app.secret
{{ end -}}