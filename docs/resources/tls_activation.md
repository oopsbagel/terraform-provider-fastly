---
layout: "fastly"
page_title: "Fastly: tls_activation"
sidebar_current: "docs-fastly-resource-tls_activation"
description: |-
Enables TLS on a domain
---

# fastly_tls_activation

Enables TLS on a domain using a specified custom TLS certificate.

~> **Note:** The Fastly service must be provisioned _prior_ to enabling TLS on it. This can be achieved in Terraform using [`depends_on`](https://www.terraform.io/docs/configuration/meta-arguments/depends_on.html).

## Example Usage

Basic usage:

```terraform
resource "fastly_service_vcl" "demo" {
  name = "my-service"

  domain {
    name = "example.com"
  }

  backend {
    address = "127.0.0.1"
    name    = "localhost"
  }

  force_destroy = true
}

resource "fastly_tls_private_key" "demo" {
  key_pem = "..."
  name    = "demo-key"
}

resource "fastly_tls_certificate" "demo" {
  certificate_body = "..."
  name             = "demo-cert"
  depends_on       = [fastly_tls_private_key.demo]
}

resource "fastly_tls_activation" "test" {
  certificate_id = fastly_tls_certificate.demo.id
  domain         = "example.com"
  depends_on     = [fastly_service_vcl.demo]
}
```

~> **Warning:** Updating the `fastly_tls_private_key`/`fastly_tls_certificate` resources should be done in multiple plan/apply steps to avoid potential downtime. The new certificate and associated private key must first be created so they exist alongside the currently active resources. Once the new resources have been created, then the `fastly_tls_activation` can be updated to point to the new certificate. Finally, the original key/certificate resources can be deleted.

## Import

A TLS activation can be imported using its ID, e.g.

```sh
$ terraform import fastly_tls_activation.demo xxxxxxxx
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **certificate_id** (String) ID of certificate to use. Must have the `domain` specified in the certificate's Subject Alternative Names.
- **domain** (String) Domain to enable TLS on. Must be assigned to an existing Fastly Service.

### Optional

- **configuration_id** (String) ID of TLS configuration to be used to terminate TLS traffic, or use the default one if missing.
- **id** (String) The ID of this resource.

### Read-Only

- **created_at** (String) Time-stamp (GMT) when TLS was enabled.
