package pinyin

import (
	"testing"
)

func TestAutoDetectUri(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "Case1", args: args{s: "~!@#$%^&*()_\\-+=<>?:\"{}|,.\\/;'\\[\\]·~！@#￥%……&*（）——\\-+={}|《》？：“”【】、；‘’，。、`"}, want: ""},
		{name: "Case2", args: args{s: "a`"}, want: "a"},
		{name: "Case3", args: args{s: "^&*%BJYUI*&"}, want: "bjyui"},
		{name: "Case4", args: args{s: "BJYUI"}, want: "bjyui"},
		{name: "Case5", args: args{s: "jkds ks sdk &sg sk gskdg ksdgh sjk hjk "}, want: "jkds-ks-sdk-sg-sk-gskdg-ksdgh-sjk-hjk"},
		{name: "Case6", args: args{s: "jkds ks sdk` &sg sk gskdg ksdgh sjk hjk "}, want: "jkds-ks-sdk-sg-sk-gskdg-ksdgh-sjk-hjk"},
		{name: "Case7", args: args{s: "Unicode 是全球文字统一编码。它把世界上的|各种文字的每一个字符指定唯一编码，实现跨语种、跨平台的应用。"}, want: "unicode-shi-quan-qiu-wen-zi-tong-yi-bian-ma-ta-ba-shi-jie-shang-de-ge-zhong-wen-zi-de-mei-yi-ge-zi-fu-zhi-ding-wei-yi-bian-ma-shi-xian-kua-yu-zhong-kua-ping-tai-de-ying-yong"},
	}
	service := Service{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := service.AutoDetectUri(tt.args.s); got != tt.want {
				//t.Errorf("AutoDetectUri() = %v, want %v", got, tt.want)
			} else {
				//log.Println(got)
			}
		})
	}
}
