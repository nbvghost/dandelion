package util

import "testing"

func TestNetworkIP(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "NetworkIP"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NetworkIP(); got == "" {
				t.Errorf("NetworkIP() = %v, 找不到本地ip地址", got)
			} else {
				t.Log(got)
			}
		})
	}
}
