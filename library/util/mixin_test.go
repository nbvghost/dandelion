package util

import (
	"reflect"
	"testing"
)

type Myst struct {
	Name string
}
func TestJSONToStruct(t *testing.T) {
	type args struct {
		j string
	}
	type testCase[T IJSON] struct {
		name    string
		args    args
		want    T
		wantErr bool
	}
	tests := []testCase[*Myst]{
		{
			name: "json#1",
			args: args{j: `{"Name":"ddd"}`},
			want: &Myst{Name: "ddd"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := JSONToStruct[*Myst](tt.args.j)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONToStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JSONToStruct() got = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestJSONToStruct1(t *testing.T) {
	type args struct {
		j string
	}
	type testCase[T IJSON] struct {
		name    string
		args    args
		want    T
		wantErr bool
	}
	tests := []testCase[[]Myst]{
		{
			name: "json#1",
			args: args{j: `[{"Name":"ddd"},{"Name":"ddd"}]`},
			want: []Myst{Myst{Name: "ddd"},Myst{Name: "ddd"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := JSONToStruct[[]Myst](tt.args.j)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONToStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JSONToStruct() got = %v, want %v", got, tt.want)
			}
		})
	}
}