# Ollama-Machine

Think `docker-machine` but for GPU cloud instances running Ollama, with a focus on safe cloud credentials storage and security.

**Disclaimer:** Currently, only two providers are available: OVHcloud and Openstack. Please note that although I work at OVHcloud, this project is not affiliated with or sponsored by OVHcloud in any way. This project is maintained on my personal free time.

## 🚀 About the Project

**Ollama-machine** is an open-source tool designed to simplify and secure the management of GPU-powered cloud instances optimized for running [Ollama](https://ollama.com). Inspired by the simplicity of `docker-machine`, this project makes it effortless to provision, manage, and secure cloud-based instances for machine learning and AI workflows.

For example, with **Ollama-machine**, you can run Ollama models on an Nvidia V100S 32GB for just $0.88/hour on OVHcloud. Work for 8 hours for only €7.04! Start your instance when you begin experimenting and remove it when you're done.

## ✨ Features

- 🛡️ **Secure Credential Management**: Your cloud provider credentials are securely stored in your system keyring, ensuring encrypted and isolated storage.
- 💻 **Streamlined Instance Setup**: Automates the provisioning of GPU-enabled cloud instances pre-configured for Ollama, including all necessary dependencies.
- 🌐 **Multi-Cloud Support**: Seamlessly switch between cloud providers like AWS, GCP, Azure, OpenStack, and more.
- 🔒 **Security-First Design**: By default, your Ollama instance is not exposed to the web. Use Tailscale or an SSH tunnel to securely connect to your instance.
- ⚙️ **Cloud Instance Management**: Easily create, list, restart, and delete GPU instances.
- 🎛️ **Flexible Configurations**: Customize instance sizes and GPU types to fit your needs.
- 💻 **Cross-Platform Compatibility**: Fully supported on macOS, Linux, and Windows.

## 🛠️ Installation

### Using pre-built binaries

You can get the latest binaries on the [release page](https://github.com/alexandrevilain/ollama-machine/releases).

### Build from source

```bash
git clone https://github.com/alexandrevilain/ollama-machine.git
cd ollama-machine
make build
```

## 🚀 Getting Started

Initialize your cloud provider:

```bash
ollama-machine credentials create dev -p openstack --openstack-identity-endpoint="https://auth.cloud.ovh.net/v3" --openstack-username="my-username" --openstack-password="my-password" --openstack-tenant-name="my-tenant-name"  --openstack-region="GRA7"
```

Create a GPU instance:

```bash
ollama-machine create my-machine --provider openstack --credentials dev --instance-type t2-le-90 --image "Debian 12 - Docker" --public
```

> [!NOTE]  
> Note the `--public` flag, asking ollama-machine to publicly expose Ollama. By default, the Ollama instance is private, and you need to run `ollama-machine tunnel [machine-name]` to get access to your instance. You can also use Tailscale if you don't want to start a tunnel.

Configure Ollama to use the instance:

```bash
eval "$(ollama-machine env my-machine)"
```

List all running instances:

```bash
ollama-machine ls
```

## 📖 Documentation

Comprehensive documentation, including examples and advanced usage, is available in the [Docs](./docs/).

## 🗓️ Roadmap

Most of the future work will be around adding new cloud and connectivity providers:

**Cloud providers:**

- [x] Openstack
- [x] OVHcloud
- [ ] Scaleway
- [ ] Linode
- [ ] AWS
- [ ] Google Cloud
- [ ] Google Cloud Run
- [ ] Azure
- [ ] DigitalOcean (when GPU will be available to everyone)
- Feel free to ask for another by raising an issue and/or submitting a Pull Request.

**Connectivity Providers:**

- [x] Tailscale
- [ ] ZeroTier
- [ ] Cloudflare Tunnel
- Feel free to ask for another by raising an issue and/or submitting a Pull Request.

## 🤝 Contributing

Contributions are welcome!

Feel free to:
- Submit issues and feature requests.
- Open pull requests for bug fixes or new features.
- Share feedback and suggestions in the discussions tab.

## 📜 License

This project is licensed under the Apache 2.0 License. See the LICENSE file for more details.

## 🙌 Acknowledgments

- Inspired by docker-machine.
- Thanks to the Ollama team for their awesome tool.
- Special thanks to the open-source community for their support and contributions.

💡 **Ready to get started?** Start provisioning secure, GPU-enabled cloud instances for Ollama today with ollama-machine!