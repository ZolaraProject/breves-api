{{- define "brevesapi.name" -}}
{{- default "user-api" .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/* Helm required labels */}}
{{- define "brevesapi.labels" -}}
heritage: {{ .Release.Service }}
release: {{ .Release.Name }}
chart: {{ .Chart.Name }}
app: "{{ template "brevesapi.name" . }}"
{{- end -}}

{{/* matchLabels */}}
{{- define "brevesapi.matchLabels" -}}
release: {{ .Release.Name }}
app: "{{ template "brevesapi.name" . }}"
{{- end -}}
