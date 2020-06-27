# netctl-obfuscate

## Description

`netctl` uses `wpa_supplicant` for connecting to wireless networks. `wpa_supplicant` comes with a tool, `wpa_passphrase`, which obfuscates the WPA/WEP key.
`netctl-obfuscate` wraps around wpa_passphrase to automatically update `netctl` configs to use the obfuscated PSK instead of the plaintext key.

## Building/Installing
### Requirements
Requires `go`.

### Building
Clone this repo (`git clone github.com/gryffyn/netctl-obfuscate/`), `cd` to it, and run `go build`. The binary will be written to the working directory.

Alternatively, run the first two commands and then `go install` which will install a binary to your `GOBIN` path.

### Installing
Run `go get github.com/gryffyn/netctl-obfuscate`.

Alternatively, if you're running Arch/Arch derivatives, `netctl-obfuscate` is available [on the AUR.](https://aur.archlinux.org/packages/netctl-obfuscate/)

## LICENSE

See the [copyright](https://github.com/gryffyn/netctl-obfuscate/blob/master/LICENSE) file.
