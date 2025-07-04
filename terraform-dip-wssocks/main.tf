locals {
  signing_key = var.signing_key != "" ? var.signing_key : random_password.signing_key.result
}

resource "random_password" "signing_key" {
  length           = 32
  special          = true
  override_special = "_%@"
}

resource "random_password" "salt" {
  length           = 16
  special          = false
  override_special = "_%@"
}

resource "random_pet" "instance" {
}

resource "hsdp_tenant_key" "key" {
  project      = "dip"
  organization = "dip"
  signing_key  = random_password.signing_key.result
  expiration   = "2025-12-31T23:59:59Z"
  salt         = random_password.salt.result
}

resource "cloudfoundry_app" "server" {
  name         = "server-${random_pet.instance.id}"
  space        = data.cloudfoundry_space.space.id
  docker_image = var.server_image
  memory       = 128
  strategy     = "blue-green"
  instances    = var.server_instances
  command      = "/app/app server"

  environment = {
    WSSOCKS_SIGNING_KEY = local.signing_key
  }

  routes {
    route = cloudfoundry_route.server.id
  }

  dynamic "routes" {
    for_each = cloudfoundry_route.pl
    content {
      route = routes.value.id
    }
  }
}

resource "cloudfoundry_route" "server" {
  domain   = data.cloudfoundry_domain.public.id
  space    = data.cloudfoundry_space.space.id
  hostname = "server-${random_pet.instance.id}"
}

resource "cloudfoundry_route" "pl" {
  count    = var.pl_host != "" ? 1 : 0
  domain   = data.cloudfoundry_domain.pl.id
  space    = data.cloudfoundry_space.space.id
  hostname = var.pl_host
}

