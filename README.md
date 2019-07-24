# octool
OpenConnect Wrapper

## Usage
```shell
go run cmd/octool/main.go
```

## Requirements
- openconnect
- ocproxy

If you use [atpons/straightforward](https://github.com/atpons/straightforward), set PATH to it.

## Settings
Configuration file path is `~/.octool.toml` is default set. If you want change, set path to env `CONFIG_FILE`.

Sample toml:
```toml
[OpenConnect]
Juniper = true                          # append opts to --juniper
User = "bob"                            # VPN Login Username
Host = "https://vpn.example.com"        # VPN Host

[OcProxy]
Port = "9090"                           # Listen this port by ocproxy and pass traffic over VPN

[StraightForward]
Enabled = true                          # Straightforward is SOCKS over SSH Proxy
Port = "9091"                           # Listen this port and pass traffic proxy over ocproxy over VPN over SSH server
User = "bob"                            # SSH Username
Host = "bastion.example.com:22"         # SSH host
IdentityFile = "/Users/bob/.ssh/id_rsa" # SSH Key
```