package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

type Signature string

const (
	EmptySignature Signature = ""
)

// PemToPrivateKey takes private key in PEM format and returns
// ecdsa.PrivateKey
func PemToPrivateKey(pemBytes []byte) (*ecdsa.PrivateKey, error) {
	// read PEM key
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	// extract ecdsa private key
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

// SignBytes signs the provided byte array using the private key provided
func SignBytes(data []byte, privateKeyPEM []byte) (Signature, error) {
	// get the private key
	privateKey, err := PemToPrivateKey(privateKeyPEM)
	if err != nil {
		return EmptySignature, err
	}

	// calculate the hasg value of the data
	hashed := sha256.Sum256(data)

	// digitally sign the hash value
	signature, err := ecdsa.SignASN1(rand.Reader, privateKey, hashed[:])
	if err != nil {
		return EmptySignature, err
	}

	// convert the signature to base64 string and return it
	base64Signature := base64.StdEncoding.EncodeToString(signature)
	return Signature(base64Signature), nil
}
