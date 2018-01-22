package vaultroletest

import (
	"testing"

	"github.com/giantswarm/vaultrole"
)

func Test_VaultRoleTest_New(t *testing.T) {
	s := New()
	_, ok := interface{}(s).(vaultrole.Interface)
	if !ok {
		t.Fatal("VaultRoleTest does not implement correct interface")
	}
}
