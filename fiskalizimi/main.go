package fiskalizimi

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fiskalizimi/proto"
	"fmt"
	"net/http"
)
import protobuf "google.golang.org/protobuf/proto"

const (
	PrivateKeyPem = `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEINOaBfsjAy90F/GCYMhkT/PibHpI5aIVxYN0YJHC7WKfoAoGCCqGSM49
AwEHoUQDQgAEjmvAipa/zaDRphq0biefLzvse7SRN3fY4SY1edOqFlzAsYv7yZ6D
nDD65d4cs918/ZMzpfA7sm/gYOFU77qHXA==
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
		Signature Signature `json:"qr_code"`
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

	// check if status code is 200 (OK)
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("response status code %d", resp.StatusCode))
	}
	return nil
}
