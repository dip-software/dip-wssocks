variable "cf_user" {
  type = string
}

variable "cf_password" {
  type = string
}

variable "cf_org_name" {
  type = string
}

variable "cf_space_name" {
  type    = string
  default = "development"
}

variable "region" {
  type    = string
  default = "us-east"
}

variable "server_image" {
  type    = string
  default = "ghcr.io/dip-software/dip-wssocks:v0.0.6"
}

variable "server_instances" {
  type    = number
  default = 2
}

variable "pl_host" {
  type        = string
  description = "PrivateLink CF host. Default is empty i.e. no mapping"
  default     = ""
}


variable "signing_key" {
  type      = string
  default   = ""
  sensitive = true
}
