package wechat

import (
	"net/http"
	"testing"
)

func TestWorkGroupBot_SendText(t *testing.T) {
	type fields struct {
		Key    string
		Client *http.Client
	}
	type args struct {
		text string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		/*{name: "TestWorkGroupBot_SendText#1",fields: fields{
			Key:    "3aa1af8a-8f1e-4cd9-b23e-b76c2265cb83",
			Client: &http.Client{},
		},args: args{text: "ddd"},wantErr: false},*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &WorkGroupBot{
				Key:    tt.fields.Key,
				Client: tt.fields.Client,
			}
			if err := m.SendText(tt.args.text); (err != nil) != tt.wantErr {
				t.Errorf("SendText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWorkGroupBot_SendMarkdown(t *testing.T) {
	type fields struct {
		Key    string
		Client *http.Client
	}
	type args struct {
		text string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		/*{name: "TestWorkGroupBot_SendMarkdown#1",fields: fields{
			Key:    "3aa1af8a-8f1e-4cd9-b23e-b76c2265cb83",
			Client: &http.Client{},
		},args: args{text: `实时新增用户反馈<font color=\"warning\">132例</font>，请相关同事注意。
         >类型:<font color=\"comment\">用户反馈</font>
         >普通用户反馈:<font color=\"comment\">117例</font>
         >VIP用户反馈:<font color=\"comment\">15例</font>`},wantErr: false},*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &WorkGroupBot{
				Key:    tt.fields.Key,
				Client: tt.fields.Client,
			}
			if err := m.SendMarkdown(tt.args.text); (err != nil) != tt.wantErr {
				t.Errorf("SendMarkdown() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}