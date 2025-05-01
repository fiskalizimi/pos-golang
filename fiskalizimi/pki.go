package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"strconv"
)

func GeneratePrivateKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

func PrivateKey2Pem(privateKey *ecdsa.PrivateKey) ([]byte, error) {
	privateBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: privateBytes}), nil
}

func CreateCSR(privateKey *ecdsa.PrivateKey, data CsrRequest) ([]byte, error) {
	subject := pkix.Name{
		Country:            []string{data.Country},
		Organization:       []string{strconv.FormatUint(data.Nui, 10)},
		OrganizationalUnit: []string{strconv.FormatUint(data.PosId, 10)},
		Locality:           []string{strconv.FormatUint(data.BranchId, 10)},
		CommonName:         data.BusinessName,
	}

	csrTemplate := x509.CertificateRequest{
		Subject:            subject,
		SignatureAlgorithm: x509.ECDSAWithSHA256,
	}

	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &csrTemplate, privateKey)
	if err != nil {
		return nil, err
	}

	pemBlock := &pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrBytes,
	}

	return pem.EncodeToMemory(pemBlock), nil
}
