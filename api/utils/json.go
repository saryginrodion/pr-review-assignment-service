package utils

import (
	"encoding/json"
	"io"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate = nil

func ValidateInstance() *validator.Validate {
	if validate == nil {
		validate = validator.New()
	}

	return validate
}

func ParseAndValidateJson[T any](b io.ReadCloser, res T) (*T, error) {
	val, err := io.ReadAll(b)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(val, &res)

	if err != nil {
		return nil, err
	}

	err = ValidateInstance().Struct(res)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

func ToJsonString[T any](data T) (string, error) {
	resBytes, err := json.Marshal(data)

	if err != nil {
		return "", err
	}

	return string(resBytes), nil
}

func FromJsonString[T any](s string) (*T, error) {
	v := new(T)
	err := json.Unmarshal([]byte(s), v)

	if err != nil {
		return nil, err
	}

	return v, nil
}

