# Ollama Machine

Welcome to the Ollama Machine documentation. This guide will help you understand the features, setup, and usage of the Ollama Machine.

## Installation

### Using pre-built binaries

You can get the latest binaries on the [release page](https://github.com/alexandrevilain/ollama-machine/releases).

### Build from source

```bash
git clone https://github.com/alexandrevilain/ollama-machine.git
cd ollama-machine
make build
```

## Creating cloud credentials

The first thing you should create is cloud credentials. These credentials will then be referenced when creating a new machine.
Each provider has its own credential flags, which you can find in the dedicated documentation page for each provider:

- [Openstack](./providers/openstack.md)
- [OVHcloud](./providers/ovhcloud.md)

For instance, with OVHcloud:

```console
$ ollama-machine credentials create dev-ovh -p ovhcloud --ovhcloud-application-key="xxx" --ovhcloud-application-secret="zzz" --ovhcloud-consumer-key="xxx" --ovhcloud-project-id="my-fake-ID"
2025/01/26 20:17:46 INFO Cloud Credentials created
```

You can list all your credentials by running the `credentials ls` command:

```console
NAME    PROVIDER
dev-ovh ovhcloud 
```

Your credentials are now created, and you're ready to create your first machine.

## Selecting the connectivity

Before creating your machine, you may have to choose a connectivity solution.

By default, `ollama-machine` uses private connectivity, which requires you to start a tunnel when using your remote Ollama instance.

You can choose to pass the `--public` flag when creating your instance, which will make Ollama listen on `0.0.0.0`. This method is the easiest but also the least secure.

You can also choose to use Tailscale to expose the Ollama server. To do so, provide the `--tailscale-auth-key` flag when creating the instance. You can get more information by reading the [Tailscale connectivity provider documentation](./connectivity/tailscale.md).

## Creating the machine

To create the machine, use the `ollama-machine create [name]` command.

You have to provide the following 2 mandatory flags:

- `--credentials` (`-c`): The name of the cloud credentials you created.
- `--provider` (`-p`): The name of the cloud provider where the machine will be created.

You may need to supply extra flags, depending on your provider:

- `--image` (`-i`): The image to use for the instance. It depends on the cloud provider. Debian 12 or Ubuntu 24 is recommended.
- `--instance-type` (`-t`): The instance type. It's the size of the VM, which may be named flavor or droplet, depending on the cloud provider.
- `--region` (`-r`): The cloud provider region where the instance will be spawned.
- `--zone` (`-z`): The zone in the region where the instance will be spawned.

For instance, with OVHcloud:

```console
$ ollama-machine create my-machine --provider ovhcloud --credentials dev-ovh --instance-type t2-le-90 --image "Debian 12" --region=GRA7
2025/01/26 20:19:40 INFO Generating SSH key pair
2025/01/26 20:19:40 INFO Generating machine config
2025/01/26 20:19:40 INFO Creating machine
2025/01/26 20:19:49 INFO Saving machine configuration to disk
2025/01/26 20:19:49 INFO Waiting for machine to be ready
2025/01/26 20:19:55 INFO Still waiting for machine to be ready
2025/01/26 20:21:08 INFO Machine ready
2025/01/26 20:21:08 INFO Waiting for Ollama to be started
2025/01/26 20:21:11 INFO Waiting for SSH to be ready err="dial tcp 135.125.89.104:22: connect: connection refused"
2025/01/26 20:21:22 INFO Still waiting for Ollama to be started err="Process exited with status 3"
2025/01/26 20:22:13 INFO Still waiting for Ollama to be started status=inactive
2025/01/26 20:22:18 INFO Ollama started
2025/01/26 20:22:18 INFO Retrieving Ollama host
2025/01/26 20:22:18 INFO Machine ready!
```

You can now configure Ollama to use the instance:

```bash
eval "$(ollama-machine env my-machine)"
```

> [!NOTE]  
> Don't forget to run `ollama-machine tunnel [machine-name]` to get access to your instance if you haven't provided the `--public` or the `--tailscale-auth-key` flags.