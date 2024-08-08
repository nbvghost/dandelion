package translate

import (
	"testing"
)

func TestCheckNotTranslate(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "#/", args: args{text: "1/3"}, want: true},
		{name: "#\\", args: args{text: "1\\3"}, want: true},
		{name: "#$21.50", args: args{text: "$21.50"}, want: true},
		{name: "#21.50", args: args{text: "21.50"}, want: true},
		{name: "#21.5kg", args: args{text: "21.5kg"}, want: true},
		{name: "#21.5Kg", args: args{text: "21.5Kg"}, want: true},
		{name: "#21.5KG", args: args{text: "21.5KG"}, want: true},
		{name: "#2024-07-24 00:46:10", args: args{text: "2024-07-24 00:46:10"}, want: true},
		{name: "#5-5", args: args{text: "5-5"}, want: true},
		{name: "#+86-13809549424", args: args{text: "+86-13809549424"}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckNotTranslate(tt.args.text); got != tt.want {
				t.Errorf("CheckNotTranslate() = %v, want %v", got, tt.want)
			}
		})
	}
}
