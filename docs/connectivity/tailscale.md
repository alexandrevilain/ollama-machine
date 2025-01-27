# Tailscale

To configure your machine with Tailscale, simply set the `--tailscale-auth-key` flag when creating your machine.

To create an auth key, follow [Tailscale's official documentation about Auth Keys](https://tailscale.com/kb/1085/auth-keys).

Then create your machine:

```
ollama-machine create my-machine --provider openstack --credentials dev --instance-type b2-7 --image "Debian 12" --tailscale-auth-key="tskey-abcdef1432341818"
```

Your Ollama instance will only be accessible through the Tailscale private IP.