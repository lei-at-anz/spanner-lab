package core

import (
	"github.com/google/uuid"
	"strings"
)

func NewRandomID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
