package main

import(
	"encoding/json"
    "fmt"
    "io"
)

func UnmarshalJSON[T any](r io.Reader, v *T) error {
    decoder := json.NewDecoder(r)
    //decoder.DisallowUnknownFields()
    if err := decoder.Decode(v); err != nil {
        return fmt.Errorf("unmarshal JSON: %w", err)
    }
    return nil
}

func MarshalJSONToString[T any](v T) (string, error) {
    jsonData, err := json.Marshal(v)
    if err != nil {
        return "", fmt.Errorf("Unable to marshal JSON to string: %w", err)
    }
    return string(jsonData), nil
}