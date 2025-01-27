# OVHcloud

## Cloud credentials

Available flag list:

| Flag name                          | Type   | Description                                       |
|------------------------------------|--------|---------------------------------------------------|
| --ovhcloud-application-key        | string | OVHcloud application key                           |
| --ovhcloud-application-secret     | string | OVHcloud application secret                        |
| --ovhcloud-consumer-key           | string | OVHcloud consumer key                              |
| --ovhcloud-endpoint               | string | OVHcloud API endpoint (default "ovh-eu")           |
| --ovhcloud-project-id             | string | OVHcloud cloud project ID (also named service-name)|

To use the OVHcloud API, you can follow the [First steps with the OVHcloud APIs](https://help.ovhcloud.com/csm/en-gb-api-getting-started-ovhcloud-api?id=kb_article_view&sysparm_article=KB0042784) tutorial.

Concretely, you have to generate these credentials via the [OVH token generation](https://api.ovh.com/createToken/?GET=/*&POST=/*&PUT=/*&DELETE=/*) page with the following rights:

- GET /cloud/project/*/image
- GET /cloud/project/*/flavor
- POST /cloud/project/*/instance
- GET /cloud/project/*/instance/*
- DELETE /cloud/project/*/instance/*

You can apply the least-privilege principle by filling the projectID in the URLs, for instance:

- GET /cloud/project/my-fake-ID/image
- GET /cloud/project/my-fake-ID/flavor
- POST /cloud/project/my-fake-ID/instance
- DELETE /cloud/project/my-fake-ID/instance/*

** How to get the project ID **

In the Public Cloud section, you can retrieve your service name ID thanks to the Copy to clipboard button.

![How to get the project ID](./assets/ovhcloud/get_service_name.png)

Exemple:

```console
ollama-machine credentials create dev-ovh -p ovhcloud --ovhcloud-application-key="xxx" --ovhcloud-application-secret="zzz" --ovhcloud-consumer-key="xxx" --ovhcloud-project-id="my-fake-ID"
```

## Creating machine

Exemple

```console
ollama-machine create my-machine --provider ovhcloud --credentials dev --instance-type t2-le-90 --image "Debian 12" --region=GRA7
```

## Using OpenStack API

You may prefer using OpenStack to start your instance, you can follow the [Openstack provider documentation](openstack.md).