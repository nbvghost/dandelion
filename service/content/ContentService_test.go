package content

import (
	"testing"

	"github.com/nbvghost/gpa/types"
)

func init() {

}
func TestContentService_FindAllContentSubType(t *testing.T) {

	type args struct {
		OID types.PrimaryKey
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "", args: args{OID: 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			/*service := ContentService{}
			if got := service.FindAllContentSubType(tt.args.OID); len(got) == 0 {
				t.Errorf("FindAllContentSubType() = %v", got)
			} else {
				t.Log(got)
			}*/
		})

	}
}
