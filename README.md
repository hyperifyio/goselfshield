# goselfshield

GoSelfShield securely embeds and deploys encrypted executables in Go, using 
passphrase-protected self-replacement for enhanced security.

## Usage

Assume you have an unencrypted executable at `./my-executable`.

With `goselfshield`, you can create an encrypted, self-executable installer:

```
goselfshield --private-key PASSPHRASE --output my-installer --source my-executable
```

After creation, you can distribute `my-installer` publicly.

To use the installer, your users need to know the passphrase used during its 
creation:

```
./my-installer --installer-private-key PASSPHRASE
```

Notes:
- Omitting `--installer-output` causes the installer to replace itself and re-execute itself after successful unpacking.
- Omitting `--installer-private-key` prompts the user to enter the passphrase from the command line.

### Full example

```
$ goselfshield -output test-executable -source tmp/gomemory --private-key df0f56ae8e54c3e028d37f36d0eda121266bf49f4ef90d34873340b3175177d0
make: Nothing to be done for `all'.
Self-installer created: test-executable
$ ./test-executable --version
INSTALLER: Please enter your private key: df0f56ae8e54c3e028d37f36d0eda121266bf49f4ef90d34873340b3175177d0
INSTALLER: Decrypted to: /Users/jhh/git/hyperifyio/goselfshield/test-executable.tmp
INSTALLER: Backup made successfully to: /Users/jhh/git/hyperifyio/goselfshield/test-executable.bak
INSTALLER: Backup removed successfully.
HG's memory game v0.0.9 by Hangover Games <info@hangover.games>
URL = https://memory.hangover.games
```