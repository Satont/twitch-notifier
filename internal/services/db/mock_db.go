package db

import (
	"context"
	"fmt"
	"github.com/satont/twitch-notifier/ent"
	"time"
)

func setupTest() (*ent.Client, error) {
	source := fmt.Sprintf("file:tests%v?mode=memory&cache=shared&_fk=1", time.Now().UnixMicro())

	entClient, err := ent.Open("sqlite3", source)
	if err != nil {
		return nil, err
	}
	if err := entClient.Schema.Create(context.Background()); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return entClient, nil
}

func teardownTest(entClient *ent.Client) {
	_ = entClient.Close()
}
