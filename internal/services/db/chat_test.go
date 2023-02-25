package db

import (
	"context"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/satont/twitch-notifier/ent"
	channel2 "github.com/satont/twitch-notifier/ent/channel"
	"github.com/satont/twitch-notifier/ent/chat"
	"testing"
)

func setupTest() (*ent.Client, error) {
	entClient, err := ent.Open("sqlite3", "file:tests?mode=memory&cache=shared&_fk=1")
	if err != nil {
		return nil, err
	}
	if err := entClient.Schema.Create(context.Background()); err != nil {
		return nil, err
	}
	return entClient, nil
}

func teardownTest(entClient *ent.Client) error {
	return entClient.Close()
}

func TestChatService_GetByID(t *testing.T) {
	entClient, err := setupTest()
	if err != nil {
		t.Fatal(err)
	}
	defer teardownTest(entClient)

	settings, err := entClient.ChatSettings.Create().Save(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	_, err = entClient.Chat.
		Create().
		SetID(uuid.New()).
		SetChatID("1").
		SetService(chat.ServiceTelegram).
		SetChatSettings(settings).
		Save(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		entClient *ent.Client
	}
	type args struct {
		chatId  string
		service chat.Service
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ent.Chat
		wantErr bool
	}{
		{
			name: "GetByID",
			fields: fields{
				entClient: entClient,
			},
			args: args{
				chatId:  "1",
				service: chat.ServiceTelegram,
			},
			want: &ent.Chat{
				ChatID:  "1",
				Service: chat.ServiceTelegram,
				Edges: ent.ChatEdges{
					ChatSettings: &ent.ChatSettings{
						ID: uuid.New(),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "GetByID",
			fields: fields{
				entClient: entClient,
			},
			args: args{
				chatId:  "2",
				service: chat.ServiceTelegram,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &chatService{
				entClient: tt.fields.entClient,
			}
			got, err := c.GetByID(tt.args.chatId, tt.args.service)
			if err != nil && !tt.wantErr {
				t.Errorf("%v error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if got != nil && got.ChatID != tt.want.ChatID {
				t.Errorf("%v got = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestChatService_Create(t *testing.T) {
	entClient, err := setupTest()
	if err != nil {
		t.Fatal(err)
	}
	defer teardownTest(entClient)

	type fields struct {
		entClient *ent.Client
	}
	type args struct {
		chatId  string
		service chat.Service
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ent.Chat
		wantErr bool
	}{
		{
			name: "Create",
			fields: fields{
				entClient: entClient,
			},
			args: args{
				chatId:  "1",
				service: chat.ServiceTelegram,
			},
			want: &ent.Chat{
				ChatID:  "1",
				Service: chat.ServiceTelegram,
				Edges: ent.ChatEdges{
					ChatSettings: &ent.ChatSettings{
						ID: uuid.New(),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Create",
			fields: fields{
				entClient: entClient,
			},
			args: args{
				chatId:  "1",
				service: chat.ServiceTelegram,
			},
			want: &ent.Chat{
				ChatID:  "1",
				Service: chat.ServiceTelegram,
				Edges: ent.ChatEdges{
					ChatSettings: &ent.ChatSettings{
						ID: uuid.New(),
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &chatService{
				entClient: tt.fields.entClient,
			}
			got, err := c.Create(tt.args.chatId, tt.args.service)
			if err != nil && !tt.wantErr {
				t.Errorf("%v error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if got != nil && got.ChatID != tt.want.ChatID {
				t.Errorf("%v got = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestChatService_GetFollowsByID(t *testing.T) {
	entClient, err := setupTest()
	if err != nil {
		t.Fatal(err)
	}
	defer teardownTest(entClient)

	settings, err := entClient.ChatSettings.Create().Save(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	channel, err := entClient.Channel.Create().SetID(uuid.New()).SetService(channel2.ServiceTwitch).SetChannelID("1").
		Save(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	ch, err := entClient.Chat.
		Create().
		SetID(uuid.New()).
		SetChatID("1").
		SetService(chat.ServiceTelegram).
		SetChatSettings(settings).
		Save(context.Background())

	_, err = entClient.Follow.Create().SetID(uuid.New()).SetChannel(channel).SetChat(ch).Save(context.Background())

	type fields struct {
		entClient *ent.Client
	}
	type args struct {
		chatId  string
		service chat.Service
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ent.FollowEdges
		wantErr bool
	}{
		{
			name: "GetFollowsByID",
			fields: fields{
				entClient: entClient,
			},
			args: args{
				chatId:  "1",
				service: chat.ServiceTelegram,
			},
			want: ent.FollowEdges{
				Channel: channel,
				Chat:    ch,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &chatService{
				entClient: tt.fields.entClient,
			}
			got, err := c.GetFollowsByID(tt.args.chatId, tt.args.service)
			if err != nil && !tt.wantErr {
				t.Errorf("%v error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if got != nil && len(got) != 1 {
				t.Errorf("%v got = %v, want %v", tt.name, got, tt.want)
			}
			if got[0].Edges.Chat.ChatID != tt.want.Chat.ChatID || got[0].Edges.Channel.ChannelID != tt.want.Channel.ChannelID {
				t.Errorf("%v got = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
