# netctl-obfuscate

![Go](https://github.com/gryffyn/netctl-obfuscate/workflows/Go/badge.svg?branch=master)

## Description

`netctl` uses `wpa_supplicant` for connecting to wireless networks. `wpa_supplicant` comes with a tool, `wpa_passphrase`, which obfuscates the WPA/WEP key.
`netctl-obfuscate` wraps around wpa_passphrase to automatically update `netctl` profiles to use the obfuscated PSK instead of the plaintext key.

Running `netctl-obfuscate` creates a backup of the netctl profile named \<profilename\>.orig, which can be restored if the program does something it's not supposed to. If that happens, opening an issue describing what happened would be much appreciated.

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
