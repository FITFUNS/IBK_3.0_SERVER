package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama-common/runtime"
)

func AfterAuthenticate(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, out *api.Session, in *api.AuthenticateCustomRequest) error {
	userId, _ := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)

	loc := time.FixedZone("KST", 9*3600)
	now := time.Now().In(loc)
	year := now.Year()
	month := int(now.Month()) // Month()는 time.Month 타입이므로 int 변환 필요
	day := now.Day()
	changeSet := ChangSet{}
	keyId := fmt.Sprintf("%d_%d_%d", year, month, day) // 문자열 결합

	objects := []*runtime.StorageRead{&runtime.StorageRead{
		Collection: "login",
		Key:        keyId,
		UserID:     userId,
	}}

	records, err := nk.StorageRead(ctx, objects)

	if err != nil {
		logger.WithField("err", err).Error("Storage read error.")
		return err
	}

	metadata := PackMessage{
		"desc": "login",
	}

	if len(records) <= 0 {
		value := PackMessage{
			"year":  int64(year),
			"month": int64(month),
			"day":   int64(day),
		}
		// value.append("info")
		jsonBytes, err := json.Marshal(value)

		if err != nil {
			logger.WithField("err", err).Error("Json Marshal error.")
			return err
		}

		jsonString := string(jsonBytes)

		nk.StorageWrite(ctx, []*runtime.StorageWrite{
			&runtime.StorageWrite{
				Collection: "login",
				Key:        keyId,
				UserID:     userId,
				Value:      jsonString}})

		changeSet["LOGIN_COUNT"] = 1
	}

	if out.Created {
		nilString := ""
		// randomName := in.GetUsername()
		randIndex := rand.Intn(len(Names)) // 올바른 랜덤 인덱스 생성
		randomName := Names[randIndex]
		nk.AccountUpdateId(ctx, userId, nilString, metadata, randomName, nilString, nilString, nilString, nilString)
		// nk.LinkDevice(ctx, userId, userId)
		changeSet[GENERAL_POINT] = 0
		changeSet[AUGUST_POINT] = 0
		changeSet[SEPTEMBER_POINT] = 0
		changeSet[OCTOBER_POINT] = 0
		changeSet[WATERSLIDE] = 0
		changeSet[GUMMY] = 0
		changeSet[MELON] = 0
		changeSet[TOTAL] = 0
		changeSet[TOTAL_AUGUST] = 0
		changeSet[TOTAL_OCTOBER] = 0
		changeSet[TOTAL_SEPTEMBER] = 0
		changeSet[TOTAL_GENERAL] = 0
	}

	nk.WalletUpdate(ctx, userId, changeSet, metadata, true)
	return nil
}
