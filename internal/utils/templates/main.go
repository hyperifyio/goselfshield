package main

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"embed"
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

//go:embed encrypted.txt
var embeddedData embed.FS

const InstallerFlagPrefix = "installer-"

func main() {

	flags := flag.NewFlagSet("Example", flag.ContinueOnError)

	buf := new(bytes.Buffer)
	flags.SetOutput(buf)

	exePathArg, err := os.Executable()
	if err != nil {
		fmt.Printf("INSTALLER: Error getting executable path: %v\n", err)
		return
	}
	exePath, err := filepath.Abs(exePathArg)
	if err != nil {
		fmt.Printf("INSTALLER: Error getting absolute path: %v\n", err)
		return
	}

	// Define flags
	outputFileArg := flags.String(InstallerFlagPrefix+"output", exePath, "set output file")
	outputFile, err := filepath.Abs(*outputFileArg)
	if err != nil {
		fmt.Printf("INSTALLER: Error getting absolute path: %v\n", err)
		return
	}

	privateKeyString := flags.String(InstallerFlagPrefix+"private-key", "", "set private key")

	isSelfInstalling := outputFile == exePath

	// Parse the flags
	err = flags.Parse(os.Args[1:])
	if err != nil {
		if strings.Contains(err.Error(), "flag provided but not defined") {
		} else {
			fmt.Printf("INSTALLER: Error: %v\n\n%s\n", err, buf.String())
			os.Exit(1)
		}
	}

	privateKey := *privateKeyString
	if privateKey == "" {

		reader := bufio.NewReader(os.Stdin)

		fmt.Print("INSTALLER: Please enter your private key: ")
		key, err := reader.ReadString('\n') // Read string up to the newline character
		if err != nil {
			fmt.Println("INSTALLER: Failed to read private key:", err)
			return
		}

		// Trim newline or carriage return depending on the platform
		key = strings.TrimSpace(key)

		privateKey = key

	}

	encryptedData, err := fs.ReadFile(embeddedData, "encrypted.txt")
	if err != nil {
		fmt.Printf("INSTALLER: Error reading embedded data: %v\n", err)
		os.Exit(1)
		return
	}

	key, err := hex.DecodeString(privateKey)
	if err != nil {
		fmt.Printf("INSTALLER: Error decoding the private key: %v\n", err)
		os.Exit(1)
		return
	}

	data, err := Decrypt(string(encryptedData), key)
	if err != nil {
		fmt.Printf("INSTALLER: Error decrypting embedded data: %v\n", err)
		os.Exit(1)
		return
	}

	tempFile := outputFile + ".tmp"
	backupFile := outputFile + ".bak"

	// Create embedded source file to tempfile
	err = os.WriteFile(tempFile, data, 0755)
	if err != nil {
		fmt.Printf("INSTALLER: Failed to unpack to file '%s': %v\n", outputFile, err)
		os.Exit(1)
		return
	}
	fmt.Println("INSTALLER: Decrypted to:", tempFile)

	// Move over the target file
	if fileExists(outputFile) {
		err = os.Rename(outputFile, backupFile)
		if err != nil {
			fmt.Printf("INSTALLER: Error moving file: %v\n", err)
			os.Exit(1)
			return
		}
		fmt.Printf("INSTALLER: Backup made successfully to: %s\n", backupFile)
	}

	// Move over the target file
	err = os.Rename(tempFile, outputFile)
	if err != nil {
		fmt.Printf("INSTALLER: Error moving file: %v\n", err)
		os.Exit(1)
		return
	}

	// Attempt to delete the backup file
	err = os.Remove(backupFile)
	if err != nil {
		// If there was an error deleting the file, print it out
		fmt.Printf("INSTALLER: Error removing file: %v\n", err)
		os.Exit(1)
		return
	}
	fmt.Println("INSTALLER: Backup removed successfully.")

	if !isSelfInstalling {
		fmt.Printf("INSTALLER: The file has been installed to: %s\n", outputFile)
		os.Exit(0)
	}

	// Prepare to execute the command
	err = syscall.Exec(exePath, filterArgs(InstallerFlagPrefix), os.Environ())
	if err != nil {
		fmt.Printf("INSTALLER: Failed to restart application: %v\n", err)
		os.Exit(1)
	}

}

// Decrypts ciphertext using AES.
func Decrypt(ciphertext string, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}

	if len(data) < gcm.NonceSize() {
		return nil, err
	}

	nonce := data[:gcm.NonceSize()]
	ciphertextBytes := data[gcm.NonceSize():]
	decryptedData, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return nil, err
	}

	return decryptedData, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func filterArgs(prefix string) []string {
	filteredArgs := []string{}
	ignoreNext := false
	for _, arg := range os.Args {
		if ignoreNext {
			ignoreNext = false
		} else {
			if strings.HasPrefix(arg, "-"+prefix) || strings.HasPrefix(arg, "--"+prefix) {
				if !strings.Contains(arg, "=") {
					ignoreNext = true
				}
			} else {
				filteredArgs = append(filteredArgs, arg)
			}
		}
	}
	return filteredArgs
}
