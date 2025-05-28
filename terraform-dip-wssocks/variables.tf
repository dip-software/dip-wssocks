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
  default = "test"
}

variable "region" {
  type    = string
  default = "us-east"
}

variable "server_image" {
  type    = string
  default = "ghcr.io/dip-software/dip-wssocks-server:v0.0.1"
}

variable "server_instances" {
  type    = number
  default = 2
}

variable "signing_key" {
  type    = string
  default = ""
  sensitive = true
}
