package baidu

import (
	"testing"
)

func TestIndex_translateBase(t *testing.T) {
	type fields struct {
		AppID  string
		AppKey string
	}
	type args struct {
		query string
		from  string
		to    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		/*{name: "#1", fields: fields{
			AppID:  "",
			AppKey: "",
		}, args: args{
			query: "这个狗是什么品种，日本种还是美国种？",
			from:  "zh",
			to:    "en",
		}, wantErr: false},*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Index{
			}
			got, err := m.translateBase(tt.args.query, tt.args.from, tt.args.to)
			if (err != nil) != tt.wantErr {
				t.Errorf("translateBase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("translateBase() got = %v", got)
		})
	}
}

func TestIndex_Translate(t *testing.T) {
	type fields struct {
		AppID  string
		AppKey string
	}
	type args struct {
		query []string
		from  string
		to    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		/*{name: "#1", fields: fields{
			AppID:  "",
			AppKey: "",
		}, args: args{
			query: []string{"这个狗是什么品种，日本种还是美国种？", "这个狗是什么品种，日本种还是美国种？", "这个狗是什么品种，日本种还是美国种？", "这个狗是什么品种，日本种还是美国种？"},
			from:  "zh",
			to:    "en",
		}, wantErr: false},
		{name: "#2", fields: fields{
			AppID:  "",
			AppKey: "",
		}, args: args{
			query: []string{
				"这个狗是什么品种，日本种还是美国种？\n中国",
				"没有高人给写稿了。准备招一些教授博导来写稿。",
				"不懂就问妈妈可以帮你破处~",
				"你的节目是依靠低人权低工资优势的",
				"没有什么比一致性更好的了。",
				"我鼓励你去你想去的地方，不要被小事所左右。\n中国",
				"专家教授还不如送外卖的大学生水平高",
				"穷且益坚啊，别坠青云之志啊",
				"我也是东方大国，凭什么中国能进，我不能进！",
				"中国说的对！",
				"我进联合国是有妙招的！",
			},
			from: "zh",
			to:   "en",
		}, wantErr: false},*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Index{

			}
			got, err := m.Translate(tt.args.query, tt.args.from, tt.args.to)
			if (err != nil) != tt.wantErr {
				t.Errorf("Translate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.args.query) {
				t.Errorf("Translate() error = %s", "翻译后，长度不一样")
				return
			}
			t.Logf("Translate() got = %v, ", got)

		})
	}
}
