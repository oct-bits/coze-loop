{{- define "application.name" -}}
    {{ printf "%s" .Chart.Name }}
{{- end -}}

{{- define "secret.name" -}}
    {{ printf "%s-secret" (include "application.name" .) }}
{{- end -}}

{{- define "image.fullname" -}}
    {{ printf "%s/%s/%s:%s" .Values.image.registry .Values.image.repository .Values.image.image .Values.image.tag }}
{{- end -}}

{{- define "configmap-etc.name" -}}
    {{ printf "%s-etc-configmap" (include "application.name" .) }}
{{- end -}}

{{- define "configmap-scripts.name" -}}
    {{ printf "%s-scripts-configmap" (include "application.name" .) }}
{{- end -}}

{{- define "configmap-init-sql.name" -}}
    {{ printf "%s-init-sql-configmap" (include "application.name" .) }}
{{- end -}}

{{- define "scripts.path" -}}
    {{ printf "/%s" (include "application.name" .) }}
{{- end -}}

{{- define "etc.path" -}}
    {{ printf "/etc/%s" (include "application.name" .) }}
{{- end -}}

{{- define "init-sql.path" -}}
    {{ printf "/etc/%s/init-sql" (include "application.name" .) }}
{{- end -}}

{{- define "entrypoint.path" -}}
    {{ printf "%s/entrypoint.sh" (include "scripts.path" .) }}
{{- end -}}

{{- define "healthcheck.path" -}}
    {{ printf "%s/healthcheck.sh" (include "scripts.path" .) }}
{{- end -}}