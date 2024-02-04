package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"

	"github.com/GDGVIT/configgy-cli/utils/config"
	"golang.org/x/crypto/pbkdf2"
)

func GeneratePassword(password string, salt []byte) (string, string, error) {
	// generate random salt for password
	if salt == nil {
		salt = make([]byte, 64)
	}
	_, err := rand.Read(salt)

	if err != nil {
		return "", "", err
	}

	key := pbkdf2.Key([]byte(password), salt, 10000, 64, sha256.New)

	keysFolderPath := filepath.Join(config.GetConfiggyCliConfigDirectory(), "keys")
	// create file to store the salt
	keyPath := filepath.Join(keysFolderPath, fmt.Sprintf("%s.key", "secretKey"))
	destinationFile, err := os.Create(keyPath)
	if err != nil {
		return "", "", err
	}
	defer destinationFile.Close()
	destinationFile.Write(salt)
	return base64.StdEncoding.EncodeToString(key), base64.RawStdEncoding.EncodeToString(salt), nil
}

func GenerateRSAKeyPair() (string, error) {
	// generate RSA key pair
	bitSize := 2048
	key, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return "", err
	}
	pub := key.Public()

	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})

	pubKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(pub.(*rsa.PublicKey)),
	})

	pubKeyPath := filepath.Join(config.GetConfiggyCliConfigDirectory(), "keys", "public.pem")
	privateKeyPath := filepath.Join(config.GetConfiggyCliConfigDirectory(), "keys", "private.pem")

	// write the keys to the file
	err = os.WriteFile(pubKeyPath, pubKeyPEM, 0644)
	if err != nil {
		return "", err
	}

	err = os.WriteFile(privateKeyPath, keyPEM, 0644)
	if err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(keyPEM), nil
}
