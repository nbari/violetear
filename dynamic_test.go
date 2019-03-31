// Test
//
// go test dynamic.go dynamic_test.go

package violetear

import "testing"

func TestSetBadName(t *testing.T) {
	s := make(dynamicSet)
	err := s.Set("test", "test")
	if err == nil {
		t.Error("Set name: test")
	}
}

func TestSetOkName(t *testing.T) {
	s := make(dynamicSet)
	err := s.Set(":test", "test")
	if err != nil {
		t.Error("Set name: :test")
	}
}

func TestRegex(t *testing.T) {
	s := make(dynamicSet)
	s.Set(":ip", `^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`)
	s.Set(":uuid", `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	uuid := "2E9C64A5-FF13-4DC5-A957-F39E39ABDC48"
	rx := s[":uuid"]
	if !rx.MatchString(uuid) {
		t.Error("regex not matching")
	}
	expect(t, len(s), 2)
}

func TestFixRegex(t *testing.T) {
	s := make(dynamicSet)
	s.Set(":name", "az")
	rx := s[":name"]
	expect(t, rx.String(), "az")
}
