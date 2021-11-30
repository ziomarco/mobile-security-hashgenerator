/*
Copyright © 2021 Marco Palmisano

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	b64 "encoding/base64"

	"github.com/spf13/cobra"
	"github.com/ziomarco/mobile-security-hashgenerator/cryptoutils"
)

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "This command can be used to decrypt a map file",
	Long:  `This command can be used to decrypt a map file.`,
	Run: func(cmd *cobra.Command, args []string) {
		var mapFilePath, _ = cmd.Flags().GetString("map-file")
		var mapFileOut, _ = cmd.Flags().GetString("map-file-out")
		var decryptionKeyFilePath, _ = cmd.Flags().GetString("key-file")
		var useB64, _ = cmd.Flags().GetBool("b64")

		if mapFilePath == "" {
			fmt.Println("map file not specified")
			return
		}

		if decryptionKeyFilePath == "" {
			fmt.Println("You must provide a decryption key")
			return
		}

		var decryptionKeyFileContent, decKeyErr = os.ReadFile(decryptionKeyFilePath)
		if decKeyErr != nil {
			panic("Error reading decryption key file")
		}

		var mapFileContent, mapErr = os.ReadFile(mapFilePath)

		if mapErr != nil {
			fmt.Println("error reading map file")
			return
		}

		var mapFile string
		var decryptionKey []byte

		mapFile = string(mapFileContent)
		decryptionKey = decryptionKeyFileContent

		if useB64 {
			decodedMapFile, _ := b64.StdEncoding.DecodeString(mapFile)
			mapFile = string(decodedMapFile)
			decryptionKey, _ = b64.StdEncoding.DecodeString(string(decryptionKey))
		}

		var decodedMapFile = cryptoutils.Decrypt(string(mapFile), decryptionKey)
		os.WriteFile(mapFileOut, []byte(decodedMapFile), 0644)
	},
}

func init() {
	rootCmd.AddCommand(decryptCmd)
	decryptCmd.Flags().String("map-file", "", "map file to decrypt")
	decryptCmd.Flags().String("key-file", "", "file containing decryption key")
	decryptCmd.Flags().String("map-file-out", "decrypted.json", "decrypted map file path (defaults to decrypted.json)")
	decryptCmd.Flags().Bool("b64", false, "User base64 for reading map and key files")
}
