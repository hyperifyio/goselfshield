# goselfshield
GoSelfShield securely embeds and deploys encrypted executables in Go, using passphrase-protected self-replacement for enhanced security.

## USAGE

Let's assume you have unencrypted `./my-executable`.

With `goselfshield` you can create encrypted self-executable installer:

`goselfshield -p PASSHPRASE -o my-installer my-executable`

Then you can publish `my-installer` as publicly available. 

To use it, your users would need to know the passphrase used when you created the self executable installer.

`./my-installer --goselfshield-p PASSHPRASE --goselfshield-o /path/to/my-executable`

* When you ommit the `--goselfshield-o` the installer would replace itself. 
* When you ommit the `--goselfshield-p` the installer would ask it from the command line.
