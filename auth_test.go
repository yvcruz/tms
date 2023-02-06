package tms

import (
	"testing"
)

// *** remplace with your credentials

func TestTodusMessageService_SendMessageToGroup(t *testing.T) {
	type fields struct {
		Config  TodusMessageServiceConfig
		Token   string
		refresh string
	}
	type args struct {
		message string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{name: testing.CoverMode(), fields: fields{
			Config: TodusMessageServiceConfig{
				Url:      "https://broadcast.mprc.cu/api/v1",
				Username: "***",
				Password: "***",
				Uid:      "***",
			},
			Token:   "***",
			refresh: "***",
		}, args: args{message: "***"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tms := TodusMessageService{
				Config:  tt.fields.Config,
				Token:   tt.fields.Token,
				refresh: tt.fields.refresh,
			}
			if got := tms.SendMessageToGroup(tt.args.message); got != tt.want {
				t.Errorf("SendMessageToGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTodusMessageService_SendMessageToUser(t *testing.T) {
	type fields struct {
		Config  TodusMessageServiceConfig
		Token   string
		refresh string
	}
	type args struct {
		username string
		message  string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{name: "sendUser", fields: fields{
			Config: TodusMessageServiceConfig{
				Url:      "https://broadcast.mprc.cu/api/v1",
				Username: "***",
				Password: "**",
			},
			Token:   "***",
			refresh: "***",
		}, args: args{
			username: "***",
			message:  "***",
		}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tms := TodusMessageService{
				Config:  tt.fields.Config,
				Token:   tt.fields.Token,
				refresh: tt.fields.refresh,
			}
			if got := tms.SendMessageToUser(tt.args.username, tt.args.message); got != tt.want {
				t.Errorf("SendMessageToUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
