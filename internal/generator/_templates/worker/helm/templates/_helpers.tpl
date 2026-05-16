{{- define "___NAME___.gomemlimit" -}}
{{- $memory := required "resources.limits.memory is required when gomemlimit.enabled is true" .memory | toString -}}
{{- $percentage := .percentage | default 70 | int -}}
{{- if or (le $percentage 0) (gt $percentage 100) -}}
{{- fail "gomemlimit.percentage must be between 1 and 100" -}}
{{- end -}}
{{- if not (regexMatch "^[0-9]+(Ki|Mi|Gi|Ti)$" $memory) -}}
{{- fail "gomemlimit supports resources.limits.memory values with Ki, Mi, Gi, or Ti suffixes" -}}
{{- end -}}
{{- $value := regexFind "^[0-9]+" $memory | int64 -}}
{{- $unit := regexFind "(Ki|Mi|Gi|Ti)$" $memory -}}
{{- $bytes := int64 0 -}}
{{- if eq $unit "Ki" -}}
{{- $bytes = mul $value 1024 -}}
{{- else if eq $unit "Mi" -}}
{{- $bytes = mul $value 1048576 -}}
{{- else if eq $unit "Gi" -}}
{{- $bytes = mul $value 1073741824 -}}
{{- else if eq $unit "Ti" -}}
{{- $bytes = mul $value 1099511627776 -}}
{{- end -}}
{{- div (mul $bytes $percentage) 100 -}}
{{- end -}}
