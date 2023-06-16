package gras

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

// Generate RSA private key and public key and save them to a file
func GenerateRSAKey(bits int) (string, string) {

	// The generatekey function uses the random data generator random to generate a pair of RSA keys with a specified number of words
	// Reader is a global, shared strong random number generator for passwords
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)

	if err != nil {
		panic(err)
	}

	// Save private key
	// Serialize the obtained Ras private key into der encoded string of ASN. 1 through x509 standard
	X509PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)

	// Save public key
	// Get data of public key
	publicKey := privateKey.PublicKey

	// X509 encoding public key
	X509PublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)

	if err != nil {
		panic(err)
	}

	// Use PEM format to encode the output of x509
	// Create a buffer to save the private key
	var privateBuffer bytes.Buffer
	var publicBuffer bytes.Buffer

	// Build a PEM. Block structure object
	privateBlock := pem.Block{Type: "RSA Private Key", Bytes: X509PrivateKey}
	publicBlock := pem.Block{Type: "RSA Public Key", Bytes: X509PublicKey}

	// Save data to buffer
	pem.Encode(&privateBuffer, &privateBlock)
	pem.Encode(&publicBuffer, &publicBlock)

	return privateBuffer.String(), publicBuffer.String()
}

// RSA encryption
func RSA_Encrypt(plainText string, pubKey string) []byte {

	//PEM decoding
	block, _ := pem.Decode([]byte(pubKey))

	//X509 decoding
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)

	if err != nil {
		panic(err)
	}

	//Type assertion
	publicKey := publicKeyInterface.(*rsa.PublicKey)

	//Encrypt plaintext
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(plainText))

	if err != nil {
		panic(err)
	}

	//Return ciphertext
	return cipherText

}

// RSA decryption
func RSA_Decrypt(cipherText string, privKey string) []byte {

	//PEM decoding
	block, _ := pem.Decode([]byte(privKey))

	//X509 decoding
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)

	if err != nil {
		panic(err)
	}

	//Decrypt the ciphertext
	plainText, _ := rsa.DecryptPKCS1v15(rand.Reader, privateKey, []byte(cipherText))

	//Return plaintext
	return plainText

}
