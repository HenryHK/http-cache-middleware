package api

import (
	"testing"
)

func TestReadAccessKey(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			"read access key",
			"ENTER YOUR ACCESS KEY HERE",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReadAccessKey(); got != tt.want {
				t.Errorf("ReadAccessKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
