package getsession

import (
	"testing"

	"github.com/bhopalg/pitwall/internal/openf1"
)

func TestGetSession(t *testing.T) {
	testcases := []struct {
		name string
	}{
		{
			name: "first test",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			openf1Clinet := openf1.Client{}

			session := New(openf1Clinet)
		})
	}
}
