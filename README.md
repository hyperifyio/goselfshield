# goselfshield

GoSelfShield securely embeds and deploys encrypted executables in Go, using 
passphrase-protected self-replacement for enhanced security.

## Usage

Assume you have an unencrypted executable at `./my-executable`.

With `goselfshield`, you can create an encrypted, self-executable installer:

```
goselfshield -p PASSPHRASE -o my-installer my-executable
```

After creation, you can distribute `my-installer` publicly.

To use the installer, your users need to know the passphrase used during its 
creation:

```
./my-installer --goselfshield-p PASSPHRASE --goselfshield-o /path/to/my-executable
```

Notes:
- Omitting `--goselfshield-o` causes the installer to replace itself.
- Omitting `--goselfshield-p` prompts the user to enter the passphrase from the command line.
