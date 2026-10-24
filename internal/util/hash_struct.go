package util

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
)

func HashStruct(data any) (string, *AppError) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", &AppError{
			Message:    "Failed to hash struct",
			HttpStatus: http.StatusInternalServerError,
			Err:        err,
		}
	}
	hash := sha256.Sum256(jsonData)
	return hex.EncodeToString(hash[:]), nil
}
