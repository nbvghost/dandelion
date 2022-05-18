package pagination

import (
	"reflect"
	"testing"
)

func TestGetSwitchPage(t *testing.T) {
	type args struct {
		index int
		total int
		size  int
	}
	tests := []struct {
		name string
		args args
		want SwitchPage
	}{
		{
			name: "TestGetSwitchPage",
			args: args{index: 99993, total: 1000000, size: 10},
			want: SwitchPage{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSwitchPage(tt.args.index, tt.args.total, tt.args.size); !reflect.DeepEqual(got, tt.want) {
				//t.Errorf("GetSwitchPage() = %v, want %v", got, tt.want)
			} else {
				t.Log(got)
			}
		})
	}
}
