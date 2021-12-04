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

// encryptCmd represents the encryptCmd command
var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "This command can be used to encrypt a map file",
	Long:  `This command can be used to encrypt a map file.`,
	Run: func(cmd *cobra.Command, args []string) {
		var key, _ = cmd.Flags().GetString("key")
		var mapFile, _ = cmd.Flags().GetString("map-file")
		var mapFileOut, _ = cmd.Flags().GetString("map-file-out")
		var encryptionKeyFile, _ = cmd.Flags().GetString("key-file-out")
		var useB64, _ = cmd.Flags().GetBool("b64")

		if mapFile == "" {
			fmt.Println("map file not specified")
			return
		}

		var mapFileContent, mapErr = os.ReadFile(mapFile)

		if mapErr != nil {
			fmt.Println("error reading map file")
			return
		}

		var cryptedMap []byte
		var encryptionKey []byte
		var keyByteArr []byte

		if key != "" {
			keyByteArr = []byte(key)
			if len(keyByteArr) != 32 {
				fmt.Println("Key must be 32 bytes long")
				return
			}
		}

		cryptedMap, encryptionKey = cryptoutils.Encrypt(string(mapFileContent), keyByteArr)

		if useB64 {
			cryptedMap = []byte(b64.StdEncoding.EncodeToString(cryptedMap))
			encryptionKey = []byte(b64.StdEncoding.EncodeToString(encryptionKey))
		}

		os.WriteFile(mapFileOut, cryptedMap, 0644)

		if key == "" {
			os.WriteFile(encryptionKeyFile, encryptionKey, 0644)
		}
	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)
	encryptCmd.Flags().String("map-file", "", "map file to encrypt")
	encryptCmd.Flags().String("key", "", "random string (32bytes length) to encrypt the map file")
	encryptCmd.Flags().String("key-file-out", "encryption.key", "file path containing encryption key (defaults to key)")
	encryptCmd.Flags().String("map-file-out", "encrypted.json", "encrypted map file path (defaults to encrypted.json)")
	encryptCmd.Flags().Bool("b64", false, "User base64 for writing map and key files")
}
