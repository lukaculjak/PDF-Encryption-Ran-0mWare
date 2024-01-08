package main

import (
	"fmt"
	"strings"
	"os"
	"path/filepath"
	"crypto/aes"
	"crypto/cipher"
)

func main() {
	folder := "pdfs"
	const aes_key = "uLY8zbyw9IVLy3fOcVZ1tikuFf2m9irP"

	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error accessing path: ", path, err)
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".enc"){
			data, err := os.ReadFile(path)
			if err != nil {
				fmt.Println("Error reading file: ", path, err)
				return err
			}

			newFileName := strings.TrimSuffix(path, ".enc")

			file, err := os.Create(newFileName)
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

			nonce, cipherText := data[:gcm.NonceSize()], data[gcm.NonceSize():]

			plainText, err := gcm.Open(nil, nonce, cipherText, nil)
			if err != nil {
				fmt.Println("Something went wrong!")
				return err
			}

			_, err = file.Write(plainText)
			if err != nil {
				fmt.Println("Something went wrong!")
				return err
			}

			err = os.Remove(path)
			if err != nil {
				fmt.Println("Something went wrong!")
				return err
			}

			fmt.Println("Decription successfull!")
		}

		return err
	})

	if err != nil {
		fmt.Println("Error walking through the files/folders", err)
		return
	}
}