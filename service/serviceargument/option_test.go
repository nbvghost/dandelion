package serviceargument

import "testing"

func BenchmarkAddAttributes(b *testing.B) {
	m := &Options{
		Attributes: make([]Option, 0),
	}
	for i := 0; i < b.N; i++ {
		m.AddAttributes(OptionsTypeAttribute, "dd", "dd", "vv")
	}
}
func TestOptions_AddAttributes(t *testing.T) {
	type fields struct {
		Attributes []Option
	}
	type args struct {
		optionsType OptionsType
		key         string
		label       string
		value       string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{name: "TestOptions_AddAttributes#1", fields: fields{Attributes: make([]Option, 0)}, args: args{optionsType: OptionsTypeAttribute, key: "dd", label: "dd", value: "vv"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Options{
				Attributes: tt.fields.Attributes,
			}
			m.AddAttributes(tt.args.optionsType, tt.args.key, tt.args.label, tt.args.value)
		})
	}
}
