{{/* The name of the application this chart installs */}}
{{- define "ack-sagemaker-controller.app.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "ack-sagemaker-controller.app.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/* The name and version as used by the chart label */}}
{{- define "ack-sagemaker-controller.chart.name-version" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/* The name of the service account to use */}}
{{- define "ack-sagemaker-controller.service-account.name" -}}
    {{ default "default" .Values.serviceAccount.name }}
{{- end -}}

{{- define "ack-sagemaker-controller.watch-namespace" -}}
{{- if eq .Values.installScope "namespace" -}}
{{ .Values.watchNamespace | default .Release.Namespace }}
{{- end -}}
{{- end -}}

{{/* The mount path for the shared credentials file */}}
{{- define "ack-sagemaker-controller.aws.credentials.secret_mount_path" -}}
{{- "/var/run/secrets/aws" -}}
{{- end -}}

{{/* The path the shared credentials file is mounted */}}
{{- define "ack-sagemaker-controller.aws.credentials.path" -}}
{{ $secret_mount_path := include "ack-sagemaker-controller.aws.credentials.secret_mount_path" . }}
{{- printf "%s/%s" $secret_mount_path .Values.aws.credentials.secretKey -}}
{{- end -}}

{{/* The rules a of ClusterRole or Role */}}
{{- define "ack-sagemaker-controller.rbac-rules" -}}
rules:
- apiGroups:
  - ""
  resources:
  - configmaps
  - secrets
  verbs:
  - get
  - list
  - patch
  - watch
- apiGroups:
  - ""
  resources:
  - namespaces
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - sagemaker.services.k8s.aws
  resources:
  - apps
  - dataqualityjobdefinitions
  - domains
  - endpointconfigs
  - endpoints
  - featuregroups
  - hyperparametertuningjobs
  - inferencecomponents
  - labelingjobs
  - modelbiasjobdefinitions
  - modelexplainabilityjobdefinitions
  - modelpackagegroups
  - modelpackages
  - modelqualityjobdefinitions
  - models
  - monitoringschedules
  - notebookinstancelifecycleconfigs
  - notebookinstances
  - pipelineexecutions
  - pipelines
  - processingjobs
  - projects
  - spaces
  - trainingjobs
  - transformjobs
  - userprofiles
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
- apiGroups:
  - sagemaker.services.k8s.aws
  resources:
  - apps/status
  - dataqualityjobdefinitions/status
  - domains/status
  - endpointconfigs/status
  - endpoints/status
  - featuregroups/status
  - hyperparametertuningjobs/status
  - inferencecomponents/status
  - labelingjobs/status
  - modelbiasjobdefinitions/status
  - modelexplainabilityjobdefinitions/status
  - modelpackagegroups/status
  - modelpackages/status
  - modelqualityjobdefinitions/status
  - models/status
  - monitoringschedules/status
  - notebookinstancelifecycleconfigs/status
  - notebookinstances/status
  - pipelineexecutions/status
  - pipelines/status
  - processingjobs/status
  - projects/status
  - spaces/status
  - trainingjobs/status
  - transformjobs/status
  - userprofiles/status
  verbs:
  - get
  - patch
  - update
- apiGroups:
  - services.k8s.aws
  resources:
  - fieldexports
  - iamroleselectors
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
- apiGroups:
  - services.k8s.aws
  resources:
  - fieldexports/status
  - iamroleselectors/status
  verbs:
  - get
  - patch
  - update
{{- end }}

{{/* Convert k/v map to string like: "key1=value1,key2=value2,..." */}}
{{- define "ack-sagemaker-controller.feature-gates" -}}
{{- $list := list -}}
{{- range $k, $v := .Values.featureGates -}}
{{- $list = append $list (printf "%s=%s" $k ( $v | toString)) -}}
{{- end -}}
{{ join "," $list }}
{{- end -}}
