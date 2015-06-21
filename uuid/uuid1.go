package uuid

import (
	"crypto/rand"
	"fmt"
	"github.com/satori/go.uuid"
)

func UUID1() uuid.UUID {
	return uuid.NewV1()
}

func UUID4() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
