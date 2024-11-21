package model

import (
	"testing"
)

func TestArray_Scan(t *testing.T) {
	type args struct {
		value interface{}
	}
	type testCase[T ITableType] struct {
		name    string
		j       Array[T]
		args    args
		wantErr bool
	}
	tests := []testCase[GoodsAttributes]{
		{
			name: "GoodsAttributes",
			j:    Array[GoodsAttributes]{},
			args: args{
				value: []byte(`[{"ID":32,"CreatedAt":"2024-04-06T03:02:55.544941+08:00","UpdatedAt":"2024-04-06T03:02:55.544941+08:00","GoodsID":118,"GroupID":10,"Name":"gfg fdsgd fsgdsf","Value":"g ds fgdsfgdsf","OID":1}, 
 {"ID":30,"CreatedAt":"2024-04-05T04:12:48.46753+08:00","UpdatedAt":"2024-04-05T04:12:48.46753+08:00","GoodsID":118,"GroupID":8,"Name":"dgh","Value":"ghfdhgdfs","OID":1}]`),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.j.Scan(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
