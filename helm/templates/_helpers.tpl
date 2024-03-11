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
  - ""
  resources:
  - secrets
  verbs:
  - get
  - list
  - patch
  - watch
- apiGroups:
  - sagemaker.services.k8s.aws
  resources:
  - apps
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
  verbs:
  - get
  - patch
  - update
- apiGroups:
  - sagemaker.services.k8s.aws
  resources:
  - dataqualityjobdefinitions
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
  - dataqualityjobdefinitions/status
  verbs:
  - get
  - patch
  - update
- apiGroups:
  - sagemaker.services.k8s.aws
  resources:
  - domains
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
  - domains/status
  verbs:
  - get
  - patch
  - update
- apiGroups:
  - sagemaker.services.k8s.aws
  resources:
  - endpointconfigs
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
  - endpointconfigs/status
  verbs:
  - get
  - patch
  - update
- apiGroups:
  - sagemaker.services.k8s.aws
  resources:
  - endpoints
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
  - endpoints/status
  verbs:
  - get
  - patch
  - update
- apiGroups:
  - sagemaker.services.k8s.aws
  resources:
  - featuregroups
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
  - featuregroups/status
  verbs:
  - get
  - patch
  - update
- apiGroups:
  - sagemaker.services.k8s.aws
  resources:
  - hyperparametertuningjobs
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
  - hyperparametertuningjobs/status
  verbs:
  - get
  - patch
  - update
- apiGroups:
  - sagemaker.services.k8s.aws
  resources:
  - modelbiasjobdefinitions
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
  - modelbiasjobdefinitions/status
  verbs:
  - get
  - patch
  - update
- apiGroups:
  - sagemaker.services.k8s.aws
  resources:
  - modelexplainabilityjobdefinitions
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
  - modelexplainabilityjobdefinitions/status
  verbs:
  - get
  - patch
  - update
- apiGroups:
  - sagemaker.services.k8s.aws
  resources:
  - modelpackagegroups
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
  - modelpackagegroups/status
  verbs:
  - get
  - patch
  - update
- apiGroups:
  - sagemaker.services.k8s.aws
  resources:
  - modelpackages
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
  - modelpackages/status
  verbs:
  - get
  - patch
  - update
- apiGroups:
  - sagemaker.services.k8s.aws
  resources:
  - modelqualityjobdefinitions
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
  - modelqualityjobdefinitions/status
  verbs:
  - get
  - patch
  - update
- apiGroups:
  - sagemaker.services.k8s.aws
  resources:
  - models
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
  - models/status
  verbs:
  - get
  - patch
  - update
- apiGroups:
  - sagemaker.services.k8s.aws
  resources:
  - monitoringschedules
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
  - monitoringschedules/status
  verbs:
  - get
  - patch
  - update
- apiGroups:
  - sagemaker.services.k8s.aws
  resources:
  - notebookinstancelifecycleconfigs
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
  - notebookinstancelifecycleconfigs/status
  verbs:
  - get
  - patch
  - update
- apiGroups:
  - sagemaker.services.k8s.aws
  resources:
  - notebookinstances
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
  - notebookinstances/status
  verbs:
  - get
  - patch
  - update
- apiGroups:
  - sagemaker.services.k8s.aws
  resources:
  - pipelineexecutions
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
  - pipelineexecutions/status
  verbs:
  - get
  - patch
  - update
- apiGroups:
  - sagemaker.services.k8s.aws
  resources:
  - pipelines
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
  - pipelines/status
  verbs:
  - get
  - patch
  - update
- apiGroups:
  - sagemaker.services.k8s.aws
  resources:
  - processingjobs
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
  - processingjobs/status
  verbs:
  - get
  - patch
  - update
- apiGroups:
  - sagemaker.services.k8s.aws
  resources:
  - trainingjobs
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
  - trainingjobs/status
  verbs:
  - get
  - patch
  - update
- apiGroups:
  - sagemaker.services.k8s.aws
  resources:
  - transformjobs
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
  - transformjobs/status
  verbs:
  - get
  - patch
  - update
- apiGroups:
  - sagemaker.services.k8s.aws
  resources:
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
  - userprofiles/status
  verbs:
  - get
  - patch
  - update
- apiGroups:
  - services.k8s.aws
  resources:
  - adoptedresources
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
  - adoptedresources/status
  verbs:
  - get
  - patch
  - update
- apiGroups:
  - services.k8s.aws
  resources:
  - fieldexports
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
  verbs:
  - get
  - patch
  - update
{{- end }}