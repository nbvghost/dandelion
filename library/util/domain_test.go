package util

import (
	"reflect"
	"testing"
)

func TestParseDomain(t *testing.T) {
	type args struct {
		host string
	}
	tests := []struct {
		name  string
		args  args
		want  []string
		want1 string
	}{
		{name: "#1", args: args{host: "bem.usokay.com"}, want: []string{"bem"}, want1: "usokay.com"},
		{name: "#2", args: args{host: "asset.sites.ink"}, want: []string{"bem"}, want1: "usokay.com"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := ParseDomain(tt.args.host)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseDomain() got = %v,%v, want %v", got, got1, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ParseDomain() got = %v,%v, want %v", got, got1, tt.want1)
			}
		})
	}
}
