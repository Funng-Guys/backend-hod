# Run ansible
This playbook as 2 external dependencies, Hetzner and Cloudflare.

You are required to have a DNS zone added in cloudflare as well as a token that can change that zone.
See Cloudflare documentation:

In order to create the server itself, you are required to have a hetzner project and a read/write token.
See Hetzner documentation:

After you have accuired the necissary tokens, run the following commands:

1. In order to create a server run:
```
ansible-playbook main.yml --tags "create" --extra-vars "cloudflare_token=<YOUR TOKEN> cloudflare_zone=example.com hcloud_token=<YOUR TOKEN>"
```

2. In order to delete a server run:
```
ansible-playbook main.yml --tags "delete" --extra-vars "cloudflare_token=<YOUR TOKEN> cloudflare_zone=example.com hcloud_token=<YOUR TOKEN>"
```
Failiure to specify a tag will result in the server being provisioned and emedietly discarded.

This will automatically provision everything required. A new server, DNS record and also automatically provision, format and mount an external hetzner volume where the persistant data will be stored. 

## Customization
You can overwrite the default variables in the `vars` section of `main.yml`. This is done by specifying the variable name along the `--extra-vars` arguments.
You can for example change the region of server via: `--extra-vars="server_location=<hetzner region>"` or server specifications with: `--extra-vars="server_type=ccx33"`
