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
	"crypto/md5"
	"encoding/hex"
	"path/filepath"
	"strings"

	"encoding/json"
	"fmt"
	"io"
	"os"

	b64 "encoding/base64"

	"github.com/spf13/cobra"
	"github.com/ziomarco/mobile-security-hashgenerator/cryptoutils"
)

type MapFile struct {
	Path string `json:"path"`
	Hash string `json:"hash"`
}

func getFileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// generateMapCmd represents the generateMap command
var generateMapCmd = &cobra.Command{
	Use:   "generate-map",
	Short: "Command to generate a map from a files list",
	Long:  "Command to generate a map from a files list",
	Run: func(cmd *cobra.Command, args []string) {
		var files []MapFile
		var parsedFilesList []string
		var exportPlainFlag, _ = cmd.Flags().GetBool("export-plain")
		var useB64, _ = cmd.Flags().GetBool("b64")
		var filesList, _ = cmd.Flags().GetStringArray("files")
		var mapFileOut, _ = cmd.Flags().GetString("map-file-out")
		var key, _ = cmd.Flags().GetString("key")
		var plainMapFileOut, _ = cmd.Flags().GetString("plain-map-file-out")
		var encryptionKeyFile, _ = cmd.Flags().GetString("key-file-out")

		for _, file := range filesList {
			err := filepath.Walk(file, func(path string, info os.FileInfo, err error) error {
				filename := strings.Replace(path, file, "", -1)
				if info.IsDir() || strings.HasPrefix(filename, ".") || strings.HasPrefix(strings.Replace(filename, "/", "", -1), ".") {
					return nil
				}
				if err != nil {
					return err
				}
				parsedFilesList = append(parsedFilesList, path)
				return nil
			})
			if err != nil {
				fmt.Println("Error reading directory:", err)
				os.Exit(1)
			}
		}

		if len(parsedFilesList) == 0 {
			fmt.Println("No files found")
			return
		}

		for _, file := range parsedFilesList {
			fmt.Println("Processing file: ", file)
			fileHash, err := getFileMD5(file)

			if err != nil {
				fmt.Println("Error getting hash for file: ", file)
				os.Exit(1)
			}

			files = append(files, MapFile{Path: file, Hash: fileHash})
		}

		var mapFileContent, _ = json.Marshal(files)

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
			cryptedMap = []byte(b64.StdEncoding.EncodeToString([]byte(cryptedMap)))
			encryptionKey = []byte(b64.StdEncoding.EncodeToString(encryptionKey))
		}

		if exportPlainFlag {
			os.WriteFile(plainMapFileOut, mapFileContent, 0644)
		}

		os.WriteFile(mapFileOut, []byte(cryptedMap), 0644)

		if key == "" {
			os.WriteFile(encryptionKeyFile, (encryptionKey), 0644)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateMapCmd)
	generateMapCmd.Flags().Bool("export-plain", false, "Export plain text map file")
	generateMapCmd.Flags().Bool("b64", false, "User base64 for writing map and key files")
	generateMapCmd.Flags().StringArray("files", []string{}, "List of files to get for generating the map, also folder supported (repeat --files arg)")
	generateMapCmd.Flags().String("key", "", "random string (32bytes length) to encrypt the map file")
	generateMapCmd.Flags().String("key-file-out", "encryption.key", "file path containing encryption key (defaults to key)")
	generateMapCmd.Flags().String("map-file-out", "encrypted.json", "encrypted map file path (defaults to encrypted.json)")
	generateMapCmd.Flags().String("plain-map-file-out", "toencrypt.json", "map file path (defaults to toencrypt.json)")
}
