package utils

import (
	"encoding/base64"
	"fmt"
	"strings"
)

func EncodeJWT(header string, payload string) string {
	return base64.StdEncoding.EncodeToString([]byte(header)) + "." + base64.StdEncoding.EncodeToString([]byte(payload))
}

func DecodeJWT(jwt string) ([]byte, []byte, error) {
	jwtArr := strings.Split(jwt, ".")

	header, err := base64.StdEncoding.DecodeString(jwtArr[0])
	if err != nil {
		return nil, nil, fmt.Errorf("header decode error")
	}

	payload, err := base64.StdEncoding.DecodeString(jwtArr[1])
	if err != nil {
		return nil, nil, fmt.Errorf("payload decode error")
	}

	return header, payload, nil
}
