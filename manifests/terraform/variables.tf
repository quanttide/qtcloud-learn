variable "region" {
  description = "阿里云地域"
  type        = string
  default     = "cn-hangzhou"
}

variable "project" {
  description = "项目名（资源命名前缀）"
  type        = string
  default     = "qtcloud-learn"
}

variable "environment" {
  description = "环境：dev / prod"
  type        = string
  default     = "prod"
}

variable "image" {
  description = "FC 容器镜像（ACR 地址）。由 CI 注入"
  type        = string
}

variable "fc_memory" {
  description = "FC 函数内存（MB）"
  type        = number
  default     = 512
}

variable "fc_timeout" {
  description = "FC 函数超时（秒）"
  type        = number
  default     = 60
}

variable "oss_access_key_id" {
  description = "OSS 持久化凭证（FC 环境变量注入）"
  type        = string
}

variable "oss_access_key_secret" {
  description = "OSS 持久化凭证（FC 环境变量注入）"
  type        = string
}
