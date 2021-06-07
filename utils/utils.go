package utils

import (
	uuid "github.com/satori/go.uuid"
)

func GenerateUUID() string {
	uid := uuid.NewV1()
	return uid.String()
}
