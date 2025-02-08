# OpenStack

## Cloud credentials

Available flag list:

| Flag name                          | Type   | Description                              |
|------------------------------------|--------|------------------------------------------|
| --openstack-domain-name            | string | OpenStack user domain name (default "Default") |
| --openstack-identity-api-version   | int    | OpenStack identity API version (default 3) |
| --openstack-identity-endpoint      | string | OpenStack identity endpoint              |
| --openstack-password               | string | OpenStack password                       |
| --openstack-password-from-stdin    | bool   | Read OpenStack password from stdin       |
| --openstack-tenant-id              | string | OpenStack tenant ID                      |
| --openstack-tenant-name            | string | OpenStack tenant name                    |
| --openstack-username               | string | OpenStack username                       |

Example:

```console
ollama-machine credentials create [credentials-name] -p openstack --openstack-identity-endpoint="https://auth.cloud.ovh.net/v3" --openstack-username="my-username" --openstack-password="my-password" --openstack-tenant-name="my-tenant-name"
```

## Creating machine

The recommended image is `Debian 12`. For instance on OVHcloud, you can create a new machine by running:

```console
ollama-machine create my-machine --provider openstack --credentials dev --instance-type t2-le-90 --image "Debian 12" --region=GRA7
```