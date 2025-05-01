package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fiskalizimi/proto"
	"fmt"
	"io"
	"net/http"
)
import protobuf "google.golang.org/protobuf/proto"

const (
	PrivateKeyPem = `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEINiuodwPKwrI8CeAJaf7IITux/c/6h1EhWTs7p1LzRfcoAoGCCqGSM49
AwEHoUQDQgAE9tA2vnPqhZu3ZtwERTWUvM4UTR4c3S0S5AOlKtVp5euUiklHB+S8
jIknSTa+If/STSuM3sA8RJEYB5iqkvQUSQ==
-----END EC PRIVATE KEY-----`
)

// SignCitizenCoupon digitally signs the CitizenCoupon provided and returns the string that
// should be encoded into QR Code
func SignCitizenCoupon(cpn *proto.CitizenCoupon) (string, error) {

	// Serialize the citizen coupon message to protobuf binary
	cpnProto, err := protobuf.Marshal(cpn)
	if err != nil {
		return "", err
	}

	// convert the serialized protobuf of citizen coupon to base64 string
	base64EncodedProto := base64.StdEncoding.EncodeToString(cpnProto)

	// digitally sign the base 64 encoded protobuf and return the signature
	signature, err := SignBytes([]byte(base64EncodedProto), []byte(PrivateKeyPem))
	if err != nil {
		return "", err
	}

	fmt.Println("Coupon    : ", base64EncodedProto)
	fmt.Println("Signature : ", signature)

	// Combine the encoded data and signature to create QR Code string and return it
	qrCodeString := fmt.Sprintf("%s|%s", base64EncodedProto, signature)
	fmt.Println("QR Code   : ", qrCodeString)

	return qrCodeString, nil
}

// SignPosCoupon digitally signs the PosCoupon provided and returns the pos coupon
// and signature in encoded base64 string
func SignPosCoupon(cpn *proto.PosCoupon) (string, Signature, error) {

	// Serialize the pos coupon message to protobuf binary
	cpnProto, err := protobuf.Marshal(cpn)
	if err != nil {
		return "", EmptySignature, err
	}

	// convert the serialized protobuf of pos coupon to base64 string
	base64EncodedProto := base64.StdEncoding.EncodeToString(cpnProto)

	// digitally sign the base 64 encoded protobuf and return the signature
	signature, err := SignBytes([]byte(base64EncodedProto), []byte(PrivateKeyPem))
	if err != nil {
		return "", EmptySignature, err
	}

	fmt.Println("Pos Coupon : ", base64EncodedProto)
	fmt.Println("Signature  : ", signature)

	// return the coupon and signature as base64 strings
	return base64EncodedProto, signature, nil
}

// SendQrCode gets CitizenCoupon, digitally signs it using provided private key and then
// sends the qr code to Fiscalization System
func SendQrCode() error {
	const url = "https://fiskalizimi.atk-ks.org/citizen/coupon"

	// get Citizen Coupon
	cpn := GetCitizenCoupon()

	// sign the citizen coupon and get QR Code
	qrCode, err := SignCitizenCoupon(cpn)
	if err != nil {
		return err
	}

	// build request body
	requestBody := struct {
		CitizenID int    `json:"citizen_id"`
		QrCode    string `json:"qr_code"`
	}{
		CitizenID: 1,
		QrCode:    qrCode,
	}

	// marshal request body to json
	jsonBody, err := json.Marshal(requestBody)
	bodyReader := bytes.NewReader(jsonBody)

	// send the post request to Fiscalization System
	resp, err := http.Post(url, "application/json", bodyReader)
	if err != nil {
		return err
	}

	defer resp.Body.Close() // Ensure the response body is closed

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.New(fmt.Sprintf("Error reading response body: %s", err.Error()))
	}

	fmt.Println("Response Body: ", string(body))

	// check if status code is 200 (OK)
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("response status code %d", resp.StatusCode))
	}
	return nil
}

// SendPosCoupon gets PosCoupon, digitally signs it using provided private key and then
// // sends the coupon details (base64 encoded) and signature to Fiscalization System
func SendPosCoupon() error {
	const url = "https://fiskalizimi.atk-ks.org/pos/coupon"

	// get POS Coupon
	cpn := GetPosCoupon()

	// sign the POS coupon and get base64 protobuf representation of coupon and signature
	cpnBase64, signature, err := SignPosCoupon(cpn)
	if err != nil {
		return err
	}

	// build request body
	requestBody := struct {
		Details   string    `json:"details"`
		Signature Signature `json:"signature"`
	}{
		Details:   cpnBase64,
		Signature: signature,
	}

	// marshal request body to json
	jsonBody, err := json.Marshal(requestBody)
	bodyReader := bytes.NewReader(jsonBody)

	resp, err := http.Post(url, "application/json", bodyReader)
	if err != nil {
		return err
	}

	defer resp.Body.Close() // Ensure the response body is closed

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.New(fmt.Sprintf("Error reading response body: %s", err.Error()))
	}

	fmt.Println("Response Body: ", string(body))

	// check if status code is 200 (OK)
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("response status code %d", resp.StatusCode))
	}
	return nil
}

func main() {
	// simulated input values
	csrRequest := CsrRequest{
		Country:      "RKS",
		BusinessName: "TEST CORP",
		Nui:          510600700,
		BranchId:     1,
		PosId:        1,
	}
	privateKey, err := GeneratePrivateKey()
	if err != nil {
		fmt.Println(err)
		return
	}

	privateKyePem, err := PrivateKey2Pem(privateKey)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Generated Private Key:")
	fmt.Println(string(privateKyePem))

	// Create CSR
	csr, err := CreateCSR(privateKey, csrRequest)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Generated CSR:")
	fmt.Println(string(csr))

	fmt.Println("Sending POS Coupon ...")
	err = SendPosCoupon()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("POS Coupon sent successfully")

	fmt.Println("-------------------------------------------------")

	fmt.Println("Sending QR Code ...")
	err = SendQrCode()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("QR Code sent successfully")
}
