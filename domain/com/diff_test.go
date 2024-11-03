package com

import (
	"reflect"
	"testing"

	"github.com/nbvghost/dandelion/entity/model"
)

func Test_diff(t *testing.T) {
	type args struct {
		from any
		to   any
	}
	tests := []struct {
		name string
		args args
		want map[string]any
	}{
		{name: "struct", args: args{
			from: &model.Admin{},
			to:   &model.Admin{Phone: "52", Name: "88"},
		}, want: map[string]any{"Phone": "52"}},
		{name: "struct", args: args{
			from: &model.Admin{},
			to:   &model.Admin{},
		}, want: map[string]any{}},
		{name: "map-1", args: args{
			from: &map[string]any{"A": 5},
			to:   &map[string]any{"A": 53},
		}, want: map[string]any{}},
		{name: "map-2", args: args{
			from: &map[string]any{"B": 5},
			to:   &map[string]any{"B": 53},
		}, want: map[string]any{"B": 53}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Diff(tt.args.from, tt.args.to, "Name", "A"); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("diff() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateStructByMap(t *testing.T) {
	type User struct {
		Name string
		Age  int
		Len  float64
	}
	type args struct {
		m          map[string]any
		structType any
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{name: "set struct", args: args{
			m:          map[string]any{"Name": "dsfds", "Age": 545, "Len": 5.5},
			structType: &User{},
		}, want: &User{
			Name: "dsfds",
			Age:  545,
			Len:  5.5,
		}},
		{name: "set struct fail", args: args{
			m:          map[string]any{"Name": "dsfds", "Age": 545, "Len": 5.6},
			structType: &User{},
		}, want: &User{
			Name: "dsfds",
			Age:  545,
			Len:  5.6,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateStructByMap(tt.args.m, tt.args.structType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateStructByMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
