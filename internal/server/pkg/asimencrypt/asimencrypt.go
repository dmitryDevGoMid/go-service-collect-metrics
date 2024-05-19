package asimencrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/config"
)

type AsimEncrypt interface {
	GenerateRsaKeyPair() (*rsa.PrivateKey, *rsa.PublicKey)
	Decrypt(ciphertext []byte) (string, error)
	GenerateKey(bits int, filename string) error
	ReadPublicKey(filename string) (*rsa.PublicKey, error)
	ReadPrivateKey(filename string) (*rsa.PrivateKey, error)
	Encrypt(pub *rsa.PublicKey, msg string) ([]byte, error)
	SetPrivateKey() error
}

type asimEncrypt struct {
	cfg            *config.Config
	pathEncryptKey string
	privateKey     *rsa.PrivateKey
}

func NewAsimEncrypt(cfg *config.Config) AsimEncrypt {
	return &asimEncrypt{cfg: cfg}
}

// Get path to keys
func (asme *asimEncrypt) GetPathToKey() string {
	pathEncryptKey := asme.cfg.PathEncrypt.PathEncryptKey

	if pathEncryptKey != "" {

		d, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
		}
		return d + "/" + pathEncryptKey
	}
	return ""

}

// Set Private key
func (asme *asimEncrypt) SetPrivateKey() error {
	//Check config by private key for decode body
	pathEncryptKey := asme.GetPathToKey()

	if pathEncryptKey != "" {
		privateKey, err := asme.ReadPrivateKey(pathEncryptKey)
		if err != nil {
			return err
		}
		asme.privateKey = privateKey
		asme.pathEncryptKey = pathEncryptKey

		asme.cfg.PathEncrypt.KeyEncryptEnbled = true
		fmt.Println(asme.privateKey)
	}

	return nil
}

// GenerateRsaKeyPair generates an RSA key pair and returns the private and public keys
func (asme *asimEncrypt) GenerateRsaKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, _ := rsa.GenerateKey(rand.Reader, 4096)
	return privkey, &privkey.PublicKey
}

// GenerateKey generates a new RSA key pair and saves the private key to a file
func (asme *asimEncrypt) GenerateKey(bits int, filename string) error {
	// Generate a new RSA key pair
	priv, public := asme.GenerateRsaKeyPair()

	// Encode the private key in PKCS#1 format
	derStream := x509.MarshalPKCS1PrivateKey(priv)

	// Write the private key to a file in PEM format
	blockPrivate := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}

	filePrivate, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer filePrivate.Close()
	err = pem.Encode(filePrivate, blockPrivate)
	if err != nil {
		return err
	}

	// Encode the public key in PKIX format
	pubkey_bytes, err := x509.MarshalPKIXPublicKey(public)
	if err != nil {
		return err
	}

	// Write the public key to a file in PEM format
	blockPublic := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubkey_bytes,
	}

	filePublic, err := os.Create("public.pem")
	if err != nil {
		return err
	}
	defer filePublic.Close()
	err = pem.Encode(filePublic, blockPublic)
	if err != nil {
		return err
	}

	fmt.Println("Private Key : ", priv)
	fmt.Println("Public key ", &priv.PublicKey)

	return nil
}

// ReadPublicKey reads a public key from a file
func (asme *asimEncrypt) ReadPublicKey(filename string) (*rsa.PublicKey, error) {
	// Read the file
	data, err := os.ReadFile(filename) // changed
	if err != nil {
		return nil, err
	}

	// Decode the file from PEM format
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, err
	}

	// Parse the public key in PKIX format
	publick, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := publick.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		break // fall through
	}
	return nil, errors.New("key type is not RSA")
}

// ReadPrivateKey reads a private key from a file
func (asme *asimEncrypt) ReadPrivateKey(filename string) (*rsa.PrivateKey, error) {
	// Read the file
	data, err := os.ReadFile(filename) // changed
	if err != nil {
		return nil, err
	}

	// Decode the file from PEM format
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, err
	}

	// Parse the private key in PKCS#1 format
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

// Encrypt encrypts a message using the public key and the OAEP method
func (asme *asimEncrypt) Encrypt(pub *rsa.PublicKey, msg string) ([]byte, error) {
	// Use an empty label (label) and the SHA-256 hash function

	label := []byte("")
	hash := sha256.New()
	hash.Write(label)
	hashedLabel := hash.Sum(nil)

	// Encrypt the message
	ciphertext, err := rsa.EncryptOAEP(
		hash,
		rand.Reader,
		pub,
		[]byte(msg),
		hashedLabel)
	if err != nil {
		return nil, err
	}

	return ciphertext, nil
}

// Decrypt decrypts an encrypted message using the private key and the OAEP method
// func (asme *asimEncrypt) Decrypt(priv *rsa.PrivateKey, ciphertext []byte) (string, error) {
func (asme *asimEncrypt) Decrypt(ciphertext []byte) (string, error) {
	// Use an empty label (label) and the SHA-256 hash function
	priv := asme.privateKey

	label := []byte("")
	hash := sha256.New()
	hash.Write(label)
	hashedLabel := hash.Sum(nil)

	// Decrypt the message
	plaintext, err := rsa.DecryptOAEP(
		hash,
		rand.Reader,
		priv,
		ciphertext,
		hashedLabel)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
