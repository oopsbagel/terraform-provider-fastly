---
layout: "fastly"
page_title: "Fastly: service_compute"
sidebar_current: "docs-fastly-resource-service-compute"
description: |-
  Provides an Fastly Compute@Edge service
---

# fastly_service_compute

Provides a Fastly Compute@Edge service. Compute@Edge is a computation platform capable of running custom binaries that you compile on your own systems and upload to Fastly. Security and portability is provided by compiling your code to [WebAssembly](https://webassembly.org/), which is run at the edge using [Lucet](https://github.com/bytecodealliance/lucet), an open-source WebAssembly runtime created by Fastly. A compute service encompasses Domains and Backends.

The Service resource requires a domain name that is correctly set up to direct
traffic to the Fastly service. See Fastly's guide on [Adding CNAME Records][fastly-cname]
on their documentation site for guidance.

## Example Usage

Basic usage:

```hcl
resource "fastly_service_compute" "demo" {
    name = "demofastly"

    domain {
      name    = "demo.notexample.com"
      comment = "demo"
    }

    backend {
      address = "127.0.0.1"
      name    = "localhost"
      port    = 80
    }

    package {
      filename = "package.tar.gz"
      source_code_hash = filesha512("package.tar.gz")
    }

    force_destroy = true
}
```



### package block

The `package` block supports uploading or modifying Wasm packages for use in a Fastly Compute@Edge service. See Fastly's documentation on
[Compute@Edge](https://www.fastly.com/products/edge-compute/serverless)


[fastly-s3]: https://docs.fastly.com/en/guides/amazon-s3
[fastly-cname]: https://docs.fastly.com/en/guides/adding-cname-records
[fastly-conditionals]: https://docs.fastly.com/en/guides/using-conditions
[fastly-sumologic]: https://developer.fastly.com/reference/api/logging/sumologic/
[fastly-gcs]: https://developer.fastly.com/reference/api/logging/gcs/

## Import

Fastly Services can be imported using their service ID, e.g.


```
$ terraform import fastly_service_compute.demo xxxxxxxxxxxxxxxxxxxx
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **domain** (Block Set, Min: 1) (see [below for nested schema](#nestedblock--domain))
- **name** (String) The unique name for the Service to create
- **package** (Block List, Min: 1, Max: 1) (see [below for nested schema](#nestedblock--package))

### Optional

- **activate** (Boolean) Conditionally prevents the Service from being activated. The apply step will continue to create a new draft version but will not activate it if this is set to `false`. Default `true`
- **backend** (Block Set) (see [below for nested schema](#nestedblock--backend))
- **bigquerylogging** (Block Set) (see [below for nested schema](#nestedblock--bigquerylogging))
- **blobstoragelogging** (Block Set) (see [below for nested schema](#nestedblock--blobstoragelogging))
- **comment** (String) Description field for the service. Default `Managed by Terraform`
- **dictionary** (Block Set) (see [below for nested schema](#nestedblock--dictionary))
- **force_destroy** (Boolean) Services that are active cannot be destroyed. In order to destroy the Service, set `force_destroy` to `true`. Default `false`
- **gcslogging** (Block Set) (see [below for nested schema](#nestedblock--gcslogging))
- **healthcheck** (Block Set) (see [below for nested schema](#nestedblock--healthcheck))
- **httpslogging** (Block Set) (see [below for nested schema](#nestedblock--httpslogging))
- **id** (String) The ID of this resource.
- **logentries** (Block Set) (see [below for nested schema](#nestedblock--logentries))
- **logging_cloudfiles** (Block Set) (see [below for nested schema](#nestedblock--logging_cloudfiles))
- **logging_datadog** (Block Set) (see [below for nested schema](#nestedblock--logging_datadog))
- **logging_digitalocean** (Block Set) (see [below for nested schema](#nestedblock--logging_digitalocean))
- **logging_elasticsearch** (Block Set) (see [below for nested schema](#nestedblock--logging_elasticsearch))
- **logging_ftp** (Block Set) (see [below for nested schema](#nestedblock--logging_ftp))
- **logging_googlepubsub** (Block Set) (see [below for nested schema](#nestedblock--logging_googlepubsub))
- **logging_heroku** (Block Set) (see [below for nested schema](#nestedblock--logging_heroku))
- **logging_honeycomb** (Block Set) (see [below for nested schema](#nestedblock--logging_honeycomb))
- **logging_kafka** (Block Set) (see [below for nested schema](#nestedblock--logging_kafka))
- **logging_kinesis** (Block Set) (see [below for nested schema](#nestedblock--logging_kinesis))
- **logging_loggly** (Block Set) (see [below for nested schema](#nestedblock--logging_loggly))
- **logging_logshuttle** (Block Set) (see [below for nested schema](#nestedblock--logging_logshuttle))
- **logging_newrelic** (Block Set) (see [below for nested schema](#nestedblock--logging_newrelic))
- **logging_openstack** (Block Set) (see [below for nested schema](#nestedblock--logging_openstack))
- **logging_scalyr** (Block Set) (see [below for nested schema](#nestedblock--logging_scalyr))
- **logging_sftp** (Block Set) (see [below for nested schema](#nestedblock--logging_sftp))
- **papertrail** (Block Set) (see [below for nested schema](#nestedblock--papertrail))
- **s3logging** (Block Set) (see [below for nested schema](#nestedblock--s3logging))
- **splunk** (Block Set) (see [below for nested schema](#nestedblock--splunk))
- **sumologic** (Block Set) (see [below for nested schema](#nestedblock--sumologic))
- **syslog** (Block Set) (see [below for nested schema](#nestedblock--syslog))
- **version_comment** (String) Description field for the version

### Read-Only

- **active_version** (Number) The currently active version of your Fastly Service
- **cloned_version** (Number) The latest cloned version by the provider. The value gets only set after running `terraform apply`

<a id="nestedblock--domain"></a>
### Nested Schema for `domain`

Required:

- **name** (String) The domain that this Service will respond to

Optional:

- **comment** (String) An optional comment about the Domain.


<a id="nestedblock--package"></a>
### Nested Schema for `package`

Required:

- **filename** (String) The path to the Wasm deployment package within your local filesystem

Optional:

- **source_code_hash** (String) Used to trigger updates. Must be set to a SHA512 hash of the package file specified with the filename. The usual way to set this is filesha512("package.tar.gz") (Terraform 0.11.12 and later) or filesha512(file("package.tar.gz")) (Terraform 0.11.11 and earlier), where "package.tar.gz" is the local filename of the Wasm deployment package


<a id="nestedblock--backend"></a>
### Nested Schema for `backend`

Required:

- **address** (String) An IPv4, hostname, or IPv6 address for the Backend
- **name** (String) Name for this Backend. Must be unique to this Service

Optional:

- **auto_loadbalance** (Boolean) Denotes if this Backend should be included in the pool of backends that requests are load balanced against. Default `true`
- **between_bytes_timeout** (Number) How long to wait between bytes in milliseconds. Default `10000`
- **connect_timeout** (Number) How long to wait for a timeout in milliseconds. Default `1000`
- **error_threshold** (Number) Number of errors to allow before the Backend is marked as down. Default `0`
- **first_byte_timeout** (Number) How long to wait for the first bytes in milliseconds. Default `15000`
- **healthcheck** (String) Name of a defined `healthcheck` to assign to this backend
- **max_conn** (Number) Maximum number of connections for this Backend. Default `200`
- **max_tls_version** (String) Maximum allowed TLS version on SSL connections to this backend.
- **min_tls_version** (String) Minimum allowed TLS version on SSL connections to this backend.
- **override_host** (String) The hostname to override the Host header
- **port** (Number) The port number on which the Backend responds. Default `80`
- **shield** (String) The POP of the shield designated to reduce inbound load. Valid values for `shield` are included in the `GET /datacenters` API response
- **ssl_ca_cert** (String) CA certificate attached to origin.
- **ssl_cert_hostname** (String) Overrides ssl_hostname, but only for cert verification. Does not affect SNI at all
- **ssl_check_cert** (Boolean) Be strict about checking SSL certs. Default `true`
- **ssl_ciphers** (String) Comma separated list of OpenSSL Ciphers to try when negotiating to the backend
- **ssl_client_cert** (String, Sensitive) Client certificate attached to origin. Used when connecting to the backend
- **ssl_client_key** (String, Sensitive) Client key attached to origin. Used when connecting to the backend
- **ssl_hostname** (String) Used for both SNI during the TLS handshake and to validate the cert
- **ssl_sni_hostname** (String) Overrides ssl_hostname, but only for SNI in the handshake. Does not affect cert validation at all
- **use_ssl** (Boolean) Whether or not to use SSL to reach the Backend. Default `false`
- **weight** (Number) The [portion of traffic](https://docs.fastly.com/en/guides/load-balancing-configuration#how-weight-affects-load-balancing) to send to this Backend. Each Backend receives weight / total of the traffic. Default `100`


<a id="nestedblock--bigquerylogging"></a>
### Nested Schema for `bigquerylogging`

Required:

- **dataset** (String) The ID of your BigQuery dataset
- **name** (String) A unique name to identify this BigQuery logging endpoint
- **project_id** (String) The ID of your GCP project
- **table** (String) The ID of your BigQuery table

Optional:

- **email** (String, Sensitive) The email for the service account with write access to your BigQuery dataset. If not provided, this will be pulled from a `FASTLY_BQ_EMAIL` environment variable
- **secret_key** (String, Sensitive) The secret key associated with the service account that has write access to your BigQuery table. If not provided, this will be pulled from the `FASTLY_BQ_SECRET_KEY` environment variable. Typical format for this is a private key in a string with newlines
- **template** (String) BigQuery table name suffix template


<a id="nestedblock--blobstoragelogging"></a>
### Nested Schema for `blobstoragelogging`

Required:

- **account_name** (String) The unique Azure Blob Storage namespace in which your data objects are stored
- **container** (String) The name of the Azure Blob Storage container in which to store logs
- **name** (String) A unique name to identify the Azure Blob Storage endpoint

Optional:

- **gzip_level** (Number) Level of Gzip compression from `0-9`. `0` means no compression. `1` is the fastest and the least compressed version, `9` is the slowest and the most compressed version. Default `0`
- **message_type** (String) How the message should be formatted. Can be either `classic`, `loggly`, `logplex` or `blank`. Default `classic`
- **path** (String) The path to upload logs to. Must end with a trailing slash. If this field is left empty, the files will be saved in the container's root path
- **period** (Number) How frequently the logs should be transferred in seconds. Default `3600`
- **public_key** (String) A PGP public key that Fastly will use to encrypt your log files before writing them to disk
- **sas_token** (String, Sensitive) The Azure shared access signature providing write access to the blob service objects. Be sure to update your token before it expires or the logging functionality will not work
- **timestamp_format** (String) `strftime` specified timestamp formatting. Default `%Y-%m-%dT%H:%M:%S.000`


<a id="nestedblock--dictionary"></a>
### Nested Schema for `dictionary`

Required:

- **name** (String) A unique name to identify this dictionary

Optional:

- **write_only** (Boolean) If `true`, the dictionary is a private dictionary, and items are not readable in the UI or via API. Default is `false`. It is important to note that changing this attribute will delete and recreate the dictionary, discard the current items in the dictionary. Using a write-only/private dictionary should only be done if the items are managed outside of Terraform

Read-Only:

- **dictionary_id** (String) The ID of the dictionary


<a id="nestedblock--gcslogging"></a>
### Nested Schema for `gcslogging`

Required:

- **bucket_name** (String) The name of the bucket in which to store the logs
- **name** (String) A unique name to identify this GCS endpoint

Optional:

- **email** (String) The email address associated with the target GCS bucket on your account. You may optionally provide this secret via an environment variable, `FASTLY_GCS_EMAIL`
- **gzip_level** (Number) Level of Gzip compression, from `0-9`. `0` is no compression. `1` is fastest and least compressed, `9` is slowest and most compressed. Default `0`
- **message_type** (String) How the message should be formatted; one of: `classic`, `loggly`, `logplex` or `blank`. Default `classic`. [Fastly Documentation](https://developer.fastly.com/reference/api/logging/gcs/)
- **path** (String) Path to store the files. Must end with a trailing slash. If this field is left empty, the files will be saved in the bucket's root path
- **period** (Number) How frequently the logs should be transferred, in seconds (Default 3600)
- **secret_key** (String, Sensitive) The secret key associated with the target gcs bucket on your account. You may optionally provide this secret via an environment variable, `FASTLY_GCS_SECRET_KEY`. A typical format for the key is PEM format, containing actual newline characters where required
- **timestamp_format** (String) specified timestamp formatting (default `%Y-%m-%dT%H:%M:%S.000`)


<a id="nestedblock--healthcheck"></a>
### Nested Schema for `healthcheck`

Required:

- **host** (String) The Host header to send for this Healthcheck
- **name** (String) A unique name to identify this Healthcheck
- **path** (String) The path to check

Optional:

- **check_interval** (Number) How often to run the Healthcheck in milliseconds. Default `5000`
- **expected_response** (Number) The status code expected from the host. Default `200`
- **http_version** (String) Whether to use version 1.0 or 1.1 HTTP. Default `1.1`
- **initial** (Number) When loading a config, the initial number of probes to be seen as OK. Default `2`
- **method** (String) Which HTTP method to use. Default `HEAD`
- **threshold** (Number) How many Healthchecks must succeed to be considered healthy. Default `3`
- **timeout** (Number) Timeout in milliseconds. Default `500`
- **window** (Number) The number of most recent Healthcheck queries to keep for this Healthcheck. Default `5`


<a id="nestedblock--httpslogging"></a>
### Nested Schema for `httpslogging`

Required:

- **name** (String) The unique name of the HTTPS logging endpoint
- **url** (String) URL that log data will be sent to. Must use the https protocol

Optional:

- **content_type** (String) Value of the `Content-Type` header sent with the request
- **header_name** (String) Custom header sent with the request
- **header_value** (String) Value of the custom header sent with the request
- **json_format** (String) Formats log entries as JSON. Can be either disabled (`0`), array of json (`1`), or newline delimited json (`2`)
- **message_type** (String) How the message should be formatted; one of: `classic`, `loggly`, `logplex` or `blank`. Default `blank`
- **method** (String) HTTP method used for request. Can be either `POST` or `PUT`. Default `POST`
- **request_max_bytes** (Number) The maximum number of bytes sent in one request
- **request_max_entries** (Number) The maximum number of logs sent in one request
- **tls_ca_cert** (String, Sensitive) A secure certificate to authenticate the server with. Must be in PEM format
- **tls_client_cert** (String, Sensitive) The client certificate used to make authenticated requests. Must be in PEM format
- **tls_client_key** (String, Sensitive) The client private key used to make authenticated requests. Must be in PEM format
- **tls_hostname** (String) Used during the TLS handshake to validate the certificate


<a id="nestedblock--logentries"></a>
### Nested Schema for `logentries`

Required:

- **name** (String) Unique name to refer to this logging setup
- **token** (String) Use token based authentication (https://logentries.com/doc/input-token/)

Optional:

- **port** (Number) The port number configured in Logentries
- **use_tls** (Boolean) Whether to use TLS for secure logging


<a id="nestedblock--logging_cloudfiles"></a>
### Nested Schema for `logging_cloudfiles`

Required:

- **access_key** (String, Sensitive) Your Cloud File account access key
- **bucket_name** (String) The name of your Cloud Files container
- **name** (String) The unique name of the Rackspace Cloud Files logging endpoint
- **user** (String) The username for your Cloud Files account

Optional:

- **gzip_level** (Number) What level of GZIP encoding to have when dumping logs (default `0`, no compression)
- **message_type** (String) How the message should be formatted. One of: `classic` (default), `loggly`, `logplex` or `blank`
- **path** (String) The path to upload logs to
- **period** (Number) How frequently log files are finalized so they can be available for reading (in seconds, default `3600`)
- **public_key** (String) The PGP public key that Fastly will use to encrypt your log files before writing them to disk
- **region** (String) The region to stream logs to. One of: DFW (Dallas), ORD (Chicago), IAD (Northern Virginia), LON (London), SYD (Sydney), HKG (Hong Kong)
- **timestamp_format** (String) The `strftime` specified timestamp formatting (default `%Y-%m-%dT%H:%M:%S.000`)


<a id="nestedblock--logging_datadog"></a>
### Nested Schema for `logging_datadog`

Required:

- **name** (String) The unique name of the Datadog logging endpoint
- **token** (String, Sensitive) The API key from your Datadog account

Optional:

- **region** (String) The region that log data will be sent to. One of `US` or `EU`. Defaults to `US` if undefined


<a id="nestedblock--logging_digitalocean"></a>
### Nested Schema for `logging_digitalocean`

Required:

- **access_key** (String, Sensitive) Your DigitalOcean Spaces account access key
- **bucket_name** (String) The name of the DigitalOcean Space
- **name** (String) The unique name of the DigitalOcean Spaces logging endpoint
- **secret_key** (String, Sensitive) Your DigitalOcean Spaces account secret key

Optional:

- **domain** (String) The domain of the DigitalOcean Spaces endpoint (default `nyc3.digitaloceanspaces.com`)
- **gzip_level** (Number) What level of Gzip encoding to have when dumping logs (default `0`, no compression)
- **message_type** (String) How the message should be formatted. One of: `classic` (default), `loggly`, `logplex` or `blank`
- **path** (String) The path to upload logs to
- **period** (Number) How frequently log files are finalized so they can be available for reading (in seconds, default `3600`)
- **public_key** (String) A PGP public key that Fastly will use to encrypt your log files before writing them to disk
- **timestamp_format** (String) `strftime` specified timestamp formatting (default `%Y-%m-%dT%H:%M:%S.000`)


<a id="nestedblock--logging_elasticsearch"></a>
### Nested Schema for `logging_elasticsearch`

Required:

- **index** (String) The name of the Elasticsearch index to send documents (logs) to
- **name** (String) The unique name of the Elasticsearch logging endpoint
- **url** (String) The Elasticsearch URL to stream logs to

Optional:

- **password** (String, Sensitive) BasicAuth password for Elasticsearch
- **pipeline** (String) The ID of the Elasticsearch ingest pipeline to apply pre-process transformations to before indexing
- **request_max_bytes** (Number) The maximum number of logs sent in one request. Defaults to `0` for unbounded
- **request_max_entries** (Number) The maximum number of bytes sent in one request. Defaults to `0` for unbounded
- **tls_ca_cert** (String, Sensitive) A secure certificate to authenticate the server with. Must be in PEM format
- **tls_client_cert** (String, Sensitive) The client certificate used to make authenticated requests. Must be in PEM format
- **tls_client_key** (String, Sensitive) The client private key used to make authenticated requests. Must be in PEM format
- **tls_hostname** (String) The hostname used to verify the server's certificate. It can either be the Common Name (CN) or a Subject Alternative Name (SAN)
- **user** (String) BasicAuth username for Elasticsearch


<a id="nestedblock--logging_ftp"></a>
### Nested Schema for `logging_ftp`

Required:

- **address** (String) The FTP address to stream logs to
- **name** (String) The unique name of the FTP logging endpoint
- **password** (String, Sensitive) The password for the server (for anonymous use an email address)
- **path** (String) The path to upload log files to. If the path ends in `/` then it is treated as a directory
- **user** (String) The username for the server (can be `anonymous`)

Optional:

- **gzip_level** (Number) Gzip Compression level. Default `0`
- **message_type** (String) How the message should be formatted (default: `classic`)
- **period** (Number) How frequently the logs should be transferred, in seconds (Default `3600`)
- **port** (Number) The port number. Default: `21`
- **public_key** (String) The PGP public key that Fastly will use to encrypt your log files before writing them to disk
- **timestamp_format** (String) specified timestamp formatting (default `%Y-%m-%dT%H:%M:%S.000`)


<a id="nestedblock--logging_googlepubsub"></a>
### Nested Schema for `logging_googlepubsub`

Required:

- **name** (String) The unique name of the Google Cloud Pub/Sub logging endpoint
- **project_id** (String) The ID of your Google Cloud Platform project
- **secret_key** (String) Your Google Cloud Platform account secret key. The `private_key` field in your service account authentication JSON
- **topic** (String) The Google Cloud Pub/Sub topic to which logs will be published
- **user** (String) Your Google Cloud Platform service account email address. The `client_email` field in your service account authentication JSON


<a id="nestedblock--logging_heroku"></a>
### Nested Schema for `logging_heroku`

Required:

- **name** (String) The unique name of the Heroku logging endpoint
- **token** (String, Sensitive) The token to use for authentication (https://www.heroku.com/docs/customer-token-authentication-token/)
- **url** (String) The URL to stream logs to


<a id="nestedblock--logging_honeycomb"></a>
### Nested Schema for `logging_honeycomb`

Required:

- **dataset** (String) The Honeycomb Dataset you want to log to
- **name** (String) The unique name of the Honeycomb logging endpoint
- **token** (String, Sensitive) The Write Key from the Account page of your Honeycomb account


<a id="nestedblock--logging_kafka"></a>
### Nested Schema for `logging_kafka`

Required:

- **brokers** (String) A comma-separated list of IP addresses or hostnames of Kafka brokers
- **name** (String) The unique name of the Kafka logging endpoint
- **topic** (String) The Kafka topic to send logs to

Optional:

- **auth_method** (String) SASL authentication method. One of: plain, scram-sha-256, scram-sha-512
- **compression_codec** (String) The codec used for compression of your logs. One of: `gzip`, `snappy`, `lz4`
- **parse_log_keyvals** (Boolean) Enables parsing of key=value tuples from the beginning of a logline, turning them into record headers
- **password** (String) SASL Pass
- **request_max_bytes** (Number) Maximum size of log batch, if non-zero. Defaults to 0 for unbounded
- **required_acks** (String) The Number of acknowledgements a leader must receive before a write is considered successful. One of: `1` (default) One server needs to respond. `0` No servers need to respond. `-1`	Wait for all in-sync replicas to respond
- **tls_ca_cert** (String, Sensitive) A secure certificate to authenticate the server with. Must be in PEM format
- **tls_client_cert** (String, Sensitive) The client certificate used to make authenticated requests. Must be in PEM format
- **tls_client_key** (String, Sensitive) The client private key used to make authenticated requests. Must be in PEM format
- **tls_hostname** (String) The hostname used to verify the server's certificate. It can either be the Common Name or a Subject Alternative Name (SAN)
- **use_tls** (Boolean) Whether to use TLS for secure logging. Can be either `true` or `false`
- **user** (String) SASL User


<a id="nestedblock--logging_kinesis"></a>
### Nested Schema for `logging_kinesis`

Required:

- **access_key** (String, Sensitive) The AWS access key to be used to write to the stream
- **name** (String) The unique name of the Kinesis logging endpoint
- **secret_key** (String, Sensitive) The AWS secret access key to authenticate with
- **topic** (String) The Kinesis stream name

Optional:

- **region** (String) The AWS region the stream resides in. (Default: `us-east-1`)


<a id="nestedblock--logging_loggly"></a>
### Nested Schema for `logging_loggly`

Required:

- **name** (String) The unique name of the Loggly logging endpoint
- **token** (String, Sensitive) The token to use for authentication (https://www.loggly.com/docs/customer-token-authentication-token/).


<a id="nestedblock--logging_logshuttle"></a>
### Nested Schema for `logging_logshuttle`

Required:

- **name** (String) The unique name of the Log Shuttle logging endpoint
- **token** (String, Sensitive) The data authentication token associated with this endpoint
- **url** (String) Your Log Shuttle endpoint URL


<a id="nestedblock--logging_newrelic"></a>
### Nested Schema for `logging_newrelic`

Required:

- **name** (String) The unique name of the New Relic logging endpoint
- **token** (String, Sensitive) The Insert API key from the Account page of your New Relic account


<a id="nestedblock--logging_openstack"></a>
### Nested Schema for `logging_openstack`

Required:

- **access_key** (String, Sensitive) Your OpenStack account access key
- **bucket_name** (String) The name of your OpenStack container
- **name** (String) The unique name of the OpenStack logging endpoint
- **url** (String) Your OpenStack auth url
- **user** (String) The username for your OpenStack account

Optional:

- **gzip_level** (Number) What level of Gzip encoding to have when dumping logs (default `0`, no compression)
- **message_type** (String) How the message should be formatted; one of: `classic`, `loggly`, `logplex` or `blank`. Default `classic`. [Fastly Documentation](https://developer.fastly.com/reference/api/logging/gcs/)
- **path** (String) Path to store the files. Must end with a trailing slash. If this field is left empty, the files will be saved in the bucket's root path
- **period** (Number) How frequently the logs should be transferred, in seconds. Default `3600`
- **public_key** (String) A PGP public key that Fastly will use to encrypt your log files before writing them to disk
- **timestamp_format** (String) specified timestamp formatting (default `%Y-%m-%dT%H:%M:%S.000`)


<a id="nestedblock--logging_scalyr"></a>
### Nested Schema for `logging_scalyr`

Required:

- **name** (String) The unique name of the Scalyr logging endpoint
- **token** (String, Sensitive) The token to use for authentication (https://www.scalyr.com/keys)

Optional:

- **region** (String) The region that log data will be sent to. One of `US` or `EU`. Defaults to `US` if undefined


<a id="nestedblock--logging_sftp"></a>
### Nested Schema for `logging_sftp`

Required:

- **address** (String) The SFTP address to stream logs to
- **name** (String) The unique name of the SFTP logging endpoint
- **path** (String) The path to upload log files to. If the path ends in `/` then it is treated as a directory
- **ssh_known_hosts** (String) A list of host keys for all hosts we can connect to over SFTP
- **user** (String) The username for the server

Optional:

- **gzip_level** (Number) What level of Gzip encoding to have when dumping logs (default `0`, no compression)
- **message_type** (String) How the message should be formatted. One of: `classic` (default), `loggly`, `logplex` or `blank`
- **password** (String, Sensitive) The password for the server. If both `password` and `secret_key` are passed, `secret_key` will be preferred
- **period** (Number) How frequently log files are finalized so they can be available for reading (in seconds, default `3600`)
- **port** (Number) The port the SFTP service listens on. (Default: `22`)
- **public_key** (String) A PGP public key that Fastly will use to encrypt your log files before writing them to disk
- **secret_key** (String, Sensitive) The SSH private key for the server. If both `password` and `secret_key` are passed, `secret_key` will be preferred
- **timestamp_format** (String) The `strftime` specified timestamp formatting (default `%Y-%m-%dT%H:%M:%S.000`)


<a id="nestedblock--papertrail"></a>
### Nested Schema for `papertrail`

Required:

- **address** (String) The address of the Papertrail endpoint
- **name** (String) A unique name to identify this Papertrail endpoint
- **port** (Number) The port associated with the address where the Papertrail endpoint can be accessed


<a id="nestedblock--s3logging"></a>
### Nested Schema for `s3logging`

Required:

- **bucket_name** (String) The name of the bucket in which to store the logs
- **name** (String) The unique name of the S3 logging endpoint

Optional:

- **domain** (String) If you created the S3 bucket outside of `us-east-1`, then specify the corresponding bucket endpoint. Example: `s3-us-west-2.amazonaws.com`
- **gzip_level** (Number) Level of Gzip compression, from `0-9`. `0` is no compression. `1` is fastest and least compressed, `9` is slowest and most compressed. Default `0`
- **message_type** (String) How the message should be formatted; one of: `classic`, `loggly`, `logplex` or `blank`. Default `classic`
- **path** (String) Path to store the files. Must end with a trailing slash. If this field is left empty, the files will be saved in the bucket's root path
- **period** (Number) How frequently the logs should be transferred, in seconds. Default `3600`
- **public_key** (String) A PGP public key that Fastly will use to encrypt your log files before writing them to disk
- **redundancy** (String) The S3 redundancy level. Should be formatted; one of: `standard`, `reduced_redundancy` or null. Default `null`
- **s3_access_key** (String, Sensitive) AWS Access Key of an account with the required permissions to post logs. It is **strongly** recommended you create a separate IAM user with permissions to only operate on this Bucket. This key will be not be encrypted. You can provide this key via an environment variable, `FASTLY_S3_ACCESS_KEY`
- **s3_secret_key** (String, Sensitive) AWS Secret Key of an account with the required permissions to post logs. It is **strongly** recommended you create a separate IAM user with permissions to only operate on this Bucket. This secret will be not be encrypted. You can provide this secret via an environment variable, `FASTLY_S3_SECRET_KEY`
- **server_side_encryption** (String) Specify what type of server side encryption should be used. Can be either `AES256` or `aws:kms`
- **server_side_encryption_kms_key_id** (String) Optional server-side KMS Key Id. Must be set if server_side_encryption is set to `aws:kms`
- **timestamp_format** (String) `strftime` specified timestamp formatting (default `%Y-%m-%dT%H:%M:%S.000`)


<a id="nestedblock--splunk"></a>
### Nested Schema for `splunk`

Required:

- **name** (String) A unique name to identify the Splunk endpoint
- **url** (String) The Splunk URL to stream logs to

Optional:

- **tls_ca_cert** (String) A secure certificate to authenticate the server with. Must be in PEM format. You can provide this certificate via an environment variable, `FASTLY_SPLUNK_CA_CERT`
- **tls_client_cert** (String) The client certificate used to make authenticated requests. Must be in PEM format.
- **tls_client_key** (String, Sensitive) The client private key used to make authenticated requests. Must be in PEM format.
- **tls_hostname** (String) The hostname used to verify the server's certificate. It can either be the Common Name or a Subject Alternative Name (SAN)
- **token** (String, Sensitive) The Splunk token to be used for authentication


<a id="nestedblock--sumologic"></a>
### Nested Schema for `sumologic`

Required:

- **name** (String) A unique name to identify this Sumologic endpoint
- **url** (String) The URL to Sumologic collector endpoint

Optional:

- **message_type** (String) How the message should be formatted; one of: `classic`, `loggly`, `logplex` or `blank`. Default `classic`. See [Fastly's Documentation on Sumologic](https://developer.fastly.com/reference/api/logging/sumologic/)


<a id="nestedblock--syslog"></a>
### Nested Schema for `syslog`

Required:

- **address** (String) A hostname or IPv4 address of the Syslog endpoint
- **name** (String) A unique name to identify this Syslog endpoint

Optional:

- **message_type** (String) How the message should be formatted; one of: `classic`, `loggly`, `logplex` or `blank`. Default `classic`
- **port** (Number) The port associated with the address where the Syslog endpoint can be accessed. Default `514`
- **tls_ca_cert** (String) A secure certificate to authenticate the server with. Must be in PEM format. You can provide this certificate via an environment variable, `FASTLY_SYSLOG_CA_CERT`
- **tls_client_cert** (String) The client certificate used to make authenticated requests. Must be in PEM format. You can provide this certificate via an environment variable, `FASTLY_SYSLOG_CLIENT_CERT`
- **tls_client_key** (String, Sensitive) The client private key used to make authenticated requests. Must be in PEM format. You can provide this key via an environment variable, `FASTLY_SYSLOG_CLIENT_KEY`
- **tls_hostname** (String) Used during the TLS handshake to validate the certificate
- **token** (String) Whether to prepend each message with a specific token
- **use_tls** (Boolean) Whether to use TLS for secure logging. Default `false`