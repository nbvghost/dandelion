package content

import (
	"testing"
)

func init() {

}
func TestContentService_FindAllContentSubType(t *testing.T) {

	type args struct {
		OID uint64
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "", args: args{OID: 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := ContentService{}
			if got := service.FindAllContentSubType(tt.args.OID); len(got) == 0 {
				t.Errorf("FindAllContentSubType() = %v", got)
			} else {
				t.Log(got)
			}
		})

	}
}
