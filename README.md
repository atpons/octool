# octool
OpenConnect Wrapper

## Usage
```shell
go run cmd/octool/main.go
```

## Settings
Configuration file path is `~/.octool.toml` is default set. If you want change, set path to env `CONFIG_FILE`.

Sample toml:
```toml
[OpenConnect]
Juniper = true
User = "bob"
Host = "https://vpn.example.com"

[OcProxy]
Port = "9090"

[StraightForward]
Enabled = true
Port = "9091"
User = "bob"
Host = "bastion.example.com:22"
IdentityFile = "/Users/bob/.ssh/id_rsa"
```