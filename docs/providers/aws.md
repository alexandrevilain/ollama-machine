# AWS

## Cloud credentials

Available flag list:

| Flag name                        | Type   | Description                              |
|----------------------------------|--------|------------------------------------------|
| aws-access-key-id                | string | AWS access key ID                        |
| aws-secret-access-key            | string | AWS secret access key                    |

Example:

```console
ollama-machine credentials create [credentials-name] -p aws --aws-access-key-id="xxx" --aws-secret-access-key="xxx"
```

## Creating machine

The recommended image is `Debian 12`, its AMI is `ami-0359cb6c0c97c6607` on `eu-west-3`. You can also provide the image name.

```console
ollama-machine create my-machine --provider aws --credentials dev-aws --instance-type g6.xlarge --image ami-0359cb6c0c97c6607 --region=eu-west-3
```
