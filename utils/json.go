package utils

import (
	"github.com/json-iterator/go"
)

func JSONEncode(data interface{}) ([]byte, error) {
	return jsoniter.Marshal(data)
}

func JSONEncodeToString(data interface{}) (string, error) {
	return jsoniter.MarshalToString(data)
}

func JSONSortEncode(data interface{}) ([]byte, error) {
	return jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(data)
}

func JSONSortEncodeToString(data interface{}) (string, error) {
	return jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(data)
}

func JSONDecode(data []byte, v interface{}) error {
	return jsoniter.Unmarshal(data, v)
}

func JSONDecodeFromString(data string, v interface{}) error {
	return jsoniter.UnmarshalFromString(data, v)
}
