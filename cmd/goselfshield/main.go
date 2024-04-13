// Copyright (c) 2024. Hangover Games <info@hangover.games>. All rights reserved.

package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"

	"goselfshield"
	"goselfshield/internal/utils"
)

func main() {

	// Define flags
	sourceFile := flag.String("s", "", "set the source file")
	outputFile := flag.String("o", "", "set the output file")
	privateKeyString := flag.String("private-key", parseStringEnv("PRIVATE_KEY", ""), "set private key")
	version := flag.Bool("version", false, "Show version information")
	initPrivateKey := flag.Bool("init-private-key", false, "Create a new private key and print it")

	// Parse the flags
	flag.Parse()

	if *version {
		fmt.Printf("%s v%s by %s\nURL = %s\n", goselfshield.Name, goselfshield.Version, goselfshield.Author, goselfshield.URL)
		return
	}

	if *initPrivateKey {
		key, err := utils.GenerateKey()
		if err != nil {
			fmt.Printf("ERROR: Failed to generate key: %v\n", err)
			os.Exit(1)
		} else {
			fmt.Printf("PRIVATE_KEY=%s\n", hex.EncodeToString(key))
		}
		return
	}

	var privateKey []byte
	if *privateKeyString == "" {
		key, err := utils.GenerateKey()
		if err != nil {
			log.Printf("ERROR: Failed to generate key: %v\n", err)
			os.Exit(1)
		} else {
			log.Printf("Initialized with a random private key: %s\n", hex.EncodeToString(key))
			privateKey = key
		}
	} else {
		key, err := hex.DecodeString(*privateKeyString)
		if err != nil {
			fmt.Printf("ERROR: Failed to decode private key: %v\n", err)
			os.Exit(1)
		} else {
			privateKey = key
		}
	}

	// Load the template source
	sourceData, err := os.ReadFile(*sourceFile)
	if err != nil {
		log.Printf("failed to read source file: %v\n", err)
		os.Exit(1)
	}

	encryptedData, err := utils.Encrypt(sourceData, privateKey)
	if err != nil {
		log.Printf("failed to encrypt source file: %v\n", err)
		os.Exit(1)
	}

	err = utils.CompileBinary([]byte(encryptedData), *outputFile)
	if err != nil {
		fmt.Printf("ERROR: Failed to compile binary: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Self-installer created: %s\n", *outputFile)

}

func parseStringEnv(key string, defaultValue string) string {
	str := os.Getenv(key)
	if str == "" {
		return defaultValue
	}
	return str
}
