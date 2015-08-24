package uuid

import (
	"crypto/rand"
	"fmt"
	"regexp"
	"testing"
)

const format = "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$"

func TestUUID(t *testing.T) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	re := regexp.MustCompile(format)
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	if !re.MatchString(uuid) {
		t.Errorf("Invalid UUID4: %s", uuid)
	}
	fmt.Println(uuid)
}
