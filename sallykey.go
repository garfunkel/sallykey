/*
 * Easy to use, cross platform tool for generating public/private SSH key pairs.
 */
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"golang.org/x/crypto/ssh"
)

const (
	PRIVATE_KEY_FILE = "id_rsa"
	PUBLIC_KEY_FILE  = "id_rsa.pub"
	BITS             = 2048
	PRE_DESCRIPTION  = "This tool can be used to generate an SSH public/private key pair." +
		"\nKeys are generated using the RSA algorithm and are 2048 bits in length." +
		"\n\nThe resulting key pair can be imported for use in programs" +
		"\nsuch as FileZilla, as well as from the command line."
	POST_DESCRIPTION = "Key pair has been generated successfully." +
		"\nPlease find the files below in your home folder:" +
		"\n\nPrivate key: %v" +
		"\nPublic key: %v"
	ERROR_DESCRIPTION = "Error generating key pair.\nPlease find error details below:\n\n"
)

func main() {
	app := app.New()
	window := app.NewWindow("SSH Key Generator")
	description := widget.NewLabel(PRE_DESCRIPTION)
	generate := widget.NewButton("Generate", nil)

	generate.OnTapped = func() {
		homeDir, err := os.UserHomeDir()

		if err != nil {
			description.SetText(ERROR_DESCRIPTION + err.Error())
		} else {
			privateKeyPath := filepath.Join(homeDir, PRIVATE_KEY_FILE)
			publicKeyPath := filepath.Join(homeDir, PUBLIC_KEY_FILE)

			if err := generateKeyPair(privateKeyPath, publicKeyPath); err == nil {
				description.SetText(fmt.Sprintf(POST_DESCRIPTION, privateKeyPath, publicKeyPath))
			} else {
				description.SetText(ERROR_DESCRIPTION + err.Error())
			}
		}

		generate.Disable()
	}

	window.SetContent(widget.NewVBox(
		description,
		widget.NewHBox(
			layout.NewSpacer(),
			generate,
			widget.NewButton("Close", func() {
				app.Quit()
			}),
			layout.NewSpacer(),
		),
	))

	window.SetFixedSize(true)
	window.CenterOnScreen()
	window.ShowAndRun()
}

func generateKeyPair(privateKeyPath, publicKeyPath string) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, BITS)

	if err != nil {
		return err
	}

	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)

	if err != nil {
		return err
	}

	privateKeyText := encodePrivateKeyToPEM(privateKey)
	publicKeyText := ssh.MarshalAuthorizedKey(publicKey)

	if ioutil.WriteFile(privateKeyPath, privateKeyText, 0600); err != nil {
		return err
	}

	return ioutil.WriteFile(publicKeyPath, publicKeyText, 0600)
}

func encodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PrivateKey(privateKey),
	})
}
