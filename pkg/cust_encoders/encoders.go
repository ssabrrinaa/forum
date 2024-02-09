package cust_encoders

import (
	"encoding/base64"
	"encoding/json"
	"net/url"
)

func EncodeParams(data error) string {
	jsonData, _ := json.Marshal(data)
	base64EncodedData := base64.StdEncoding.EncodeToString(jsonData)
	params := url.Values{}
	params.Set("params", base64EncodedData)
	encodedParams := params.Encode()
	return encodedParams
}

func DecodeParams(params string) error {
	decodedParams, err := url.QueryUnescape(params)
	if err != nil {
		return err
	}

	queryParams, err := url.ParseQuery(decodedParams)
	if err != nil {
		return err
	}

	base64EncodedData := queryParams.Get("params")

	jsonData, err := base64.StdEncoding.DecodeString(base64EncodedData)
	if err != nil {
		return err
	}

	var errorData error
	err = json.Unmarshal(jsonData, &errorData)
	if err != nil {
		return err
	}

	return errorData
}
