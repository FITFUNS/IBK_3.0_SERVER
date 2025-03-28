package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/heroiclabs/nakama-common/runtime"
)

func GetMatchId(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	limit := 10
	authoritative := true
	label := payload
	minSize := 0
	maxSize := 30

	matches, err := nk.MatchList(ctx, limit, authoritative, label, &minSize, &maxSize, "")

	if err != nil {
		logger.Error("error : %v", err)

		return "", err
	}

	if len(matches) > 0 {
		return matches[0].MatchId, nil
	}

	matchId, err := nk.MatchCreate(ctx, payload, nil)

	if err != nil {
		logger.Error("error : %v", err)

		return "", err
	}

	return matchId, nil
}

func RemoveAccount(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	userID, _ := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)

	nk.UsersBanId(ctx, []string{userID})
	return "true", nil
}

func GetRank(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	userID, _ := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)

	targetId := ""
	if payload == "0" {
		targetId = GENERAL_LEADERBOARD
	} else if payload == "1" {
		targetId = OCTOBER_LEADERBOARD

	} else if payload == "2" {
		targetId = SEPTEMBER_LEADERBOARD
	} else {
		targetId = AUGUST_LEADERBOARD
	}

	records, ownerRecords, _, _, err := nk.LeaderboardRecordsList(ctx, targetId, []string{userID}, 50, "", int64(0))

	if err != nil {
		return "", err
	}

	recordsJson, err := json.Marshal([]interface{}{
		records,
		ownerRecords,
	})

	if err != nil {
		return "", err
	}

	return string(recordsJson), nil
}

func SubmitGummyRank(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	userID, _ := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)

	data := PackMessage{}

	err := json.Unmarshal([]byte(payload), &data)

	if err != nil {
		logger.Error("error : %v", err)

		return "", err
	}

	account, err := nk.AccountGetId(ctx, userID)

	if err != nil {
		logger.Error("error : %v", err)

		return "", err
	}

	metadataJson := account.User.Metadata
	metadata := PackMessage{}
	err = json.Unmarshal([]byte(metadataJson), &metadata)

	if err != nil {
		logger.Error("error : %v", err)

		return "", err
	}

	score := int64(0)

	decodedScore := data["score"].(float64)
	if decodedScore < 1 {
		score = 550
	} else if decodedScore <= 20 {
		score = 450
	} else if decodedScore <= 40 {
		score = 350
	} else if decodedScore <= 60 {
		score = 250
	} else if decodedScore <= 80 {
		score = 120
	} else if decodedScore <= 90 {
		score = 50
	}

	updated, _, _ := nk.WalletUpdate(ctx, userID, ChangSet{
		GENERAL_POINT: score,
		GUMMY:         1,
		TOTAL:         1,
		TOTAL_GENERAL: 1,
	}, PackMessage{"desc": GUMMY}, true)

	nk.LeaderboardRecordWrite(ctx, GENERAL_LEADERBOARD, userID,
		account.User.DisplayName+"#"+account.User.Username, updated[GENERAL_POINT], 0, nil, nil)

	updated["earn"] = score
	updated["best"] = metadata[GUMMY_BEST].(int64)

	res, err := json.Marshal(updated)
	if err != nil {
		logger.Error("error : %v", err)

		return "", err
	}

	return string(res), nil
}

func SubmitWaterRank(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	userID, _ := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)

	data := PackMessage{}

	err := json.Unmarshal([]byte(payload), &data)

	if err != nil {
		logger.Error("error : %v", err)

		return "", err
	}

	account, err := nk.AccountGetId(ctx, userID)

	if err != nil {
		logger.Error("error : %v", err)

		return "", err
	}

	metadataJson := account.User.Metadata
	metadata := PackMessage{}
	err = json.Unmarshal([]byte(metadataJson), &metadata)

	if err != nil {
		logger.Error("error : %v", err)

		return "", err
	}

	score := int64(0)

	decodedScore := data["score"].(float64)
	if decodedScore <= 20 {
		score = 60
	} else if decodedScore <= 40 {
		score = 80
	} else if decodedScore <= 60 {
		score = 100
	} else if decodedScore <= 80 {
		score = 150
	} else if decodedScore <= 100 {
		score = 200
	} else if decodedScore <= 120 {
		score = 300
	} else if decodedScore <= 149 {
		score = 400
	} else {
		score = 1000
	}

	updated, _, _ := nk.WalletUpdate(ctx, userID, ChangSet{
		GENERAL_POINT: score,
		WATERSLIDE:    1,
		TOTAL:         1,
		TOTAL_GENERAL: 1,
	}, PackMessage{"desc": WATERSLIDE}, true)

	nk.LeaderboardRecordWrite(ctx, GENERAL_LEADERBOARD, userID,
		account.User.DisplayName+"#"+account.User.Username, updated[GENERAL_POINT], 0, nil, nil)

	updated["earn"] = score
	updated["best"] = metadata[WATER_BEST].(int64)

	res, err := json.Marshal(updated)
	if err != nil {
		logger.Error("error : %v", err)

		return "", err
	}

	return string(res), nil
}

func SubmitMelonRank(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	userID, _ := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)

	data := PackMessage{}

	err := json.Unmarshal([]byte(payload), &data)

	if err != nil {
		logger.Error("error : %v", err)

		return "", err
	}

	account, err := nk.AccountGetId(ctx, userID)

	if err != nil {
		logger.Error("error : %v", err)

		return "", err
	}

	metadataJson := account.User.Metadata
	metadata := PackMessage{}
	err = json.Unmarshal([]byte(metadataJson), &metadata)

	if err != nil {
		logger.Error("error : %v", err)

		return "", err
	}

	decodedScore := int64(data["score"].(float64))

	updated, _, _ := nk.WalletUpdate(ctx, userID, ChangSet{
		GENERAL_POINT: decodedScore,
		MELON:         1,
		TOTAL:         1,
		TOTAL_GENERAL: 1,
	}, PackMessage{"desc": MELON}, true)

	nk.LeaderboardRecordWrite(ctx, GENERAL_LEADERBOARD, userID,
		account.User.DisplayName+"#"+account.User.Username, updated[GENERAL_POINT], 0, nil, nil)

	updated["earn"] = decodedScore
	updated["best"] = metadata[WATER_BEST].(int64)

	res, err := json.Marshal(updated)
	if err != nil {
		logger.Error("error : %v", err)

		return "", err
	}

	return string(res), nil
}

func GetWallet(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	userID, _ := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)

	account, err := nk.AccountGetId(ctx, userID)

	if err != nil {
		logger.Error("error : %v", err)

		return "", err
	}

	return account.Wallet, nil
}

func GetFitfunsReward(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	userID, _ := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)

	account, err := nk.AccountGetId(ctx, userID)

	if err != nil {
		logger.Error("error : %v", err)

		return "", err
	}

	metadataJson := account.User.Metadata
	metadata := PackMessage{}
	err = json.Unmarshal([]byte(metadataJson), &metadata)

	if err != nil {
		logger.Error("error : %v", err)

		return "", err
	}

	loc := time.FixedZone("KST", 9*3600)
	now := time.Now().In(loc)
	month := int(now.Month())
	day := now.Day()

	lastReward, ok := metadata["last_reward"].([]int)
	updateReward := []int{month, day}
	success := false

	if !ok || lastReward == nil {
		metadata["last_reward"] = updateReward
		success = true
	} else if lastReward[0] != updateReward[0] || lastReward[1] != updateReward[1] {
		metadata["last_reward"] = updateReward
		success = true
	}

	if success {
		nk.AccountUpdateId(ctx, userID, "", metadata, "", "", "", "", "")

		updated, _, _ := nk.WalletUpdate(ctx, userID, ChangSet{
			GENERAL_POINT: 200,
		}, PackMessage{"desc": "FITFUNS_REWARD"}, true)

		nk.LeaderboardRecordWrite(ctx, GENERAL_LEADERBOARD, userID,
			account.User.DisplayName+"#"+account.User.Username, updated[GENERAL_POINT], 0, nil, nil)

		res, _ := json.Marshal(PackMessage{
			"reward": 0,
		})

		return string(res), nil
	}

	res, _ := json.Marshal(PackMessage{
		"reward": 0,
	})

	return string(res), nil
}

func GetCalinder(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	userID, _ := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)

	storageData, _, err := nk.StorageList(ctx, "", userID, "login", 100, "")

	if err != nil {
		logger.Error("error : %v", err)

		return "", err
	}

	data := []interface{}{}

	for _, v := range storageData {
		value := PackMessage{}

		json.Unmarshal([]byte(v.Value), &value)
		data = append(data, value)
	}

	res, err := json.Marshal(data)

	if err != nil {
		logger.Error("error : %v", err)

		return "", err
	}

	return string(res), nil
}

func UpdatePrivacy(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	userID, _ := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)

	data := PackMessage{}
	err := json.Unmarshal([]byte(payload), &data)

	if err != nil {
		logger.WithField("err", err).Error("Unmarshal error.")
		return "", err
	}

	value, err := json.Marshal(PackMessage{
		"name":  data["name"],
		"birth": data["birth"],
		"phone": data["phone"],
	})

	if err != nil {
		logger.WithField("err", err).Error("Marshal error.")
		return "", err
	}

	_, err = nk.StorageWrite(ctx, []*runtime.StorageWrite{
		{
			Collection: "privacy",
			Key:        "privacy",
			UserID:     userID,
			Value:      string(value),
		},
	})

	if err != nil {
		logger.WithField("err", err).Error("StorageWrite error.")
		return "", err
	}

	res, _ := json.Marshal(PackMessage{
		"result": true,
	})

	return string(res), nil
}

func GetPrivacy(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	userID, _ := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)

	storageData, err := nk.StorageRead(ctx, []*runtime.StorageRead{
		{
			Collection: "privacy",
			Key:        "privacy",
			UserID:     userID,
		},
	})

	if err != nil {
		logger.WithField("err", err).Error("StorageRead error.")
		return "", err
	}

	if len(storageData) <= 0 {
		return "", runtime.NewError("Empty Privacy", 13)
	}

	return storageData[0].Value, nil
}

func UpdateAllow(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	userID, _ := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)

	account, err := nk.AccountGetId(ctx, userID)

	if err != nil {
		logger.Error("error : %v", err)

		return "", err
	}

	metadataJson := account.User.Metadata
	metadata := PackMessage{}
	err = json.Unmarshal([]byte(metadataJson), &metadata)

	if err != nil {
		logger.Error("error : %v", err)

		return "", err
	}

	metadata["is_allow"] = true

	nk.AccountUpdateId(ctx, userID, "", metadata, "", "", "", "", "")

	return "", nil
}

func GetPocketarkId(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	data := PackMessage{}

	err := json.Unmarshal([]byte(payload), &data)

	if err != nil {
		logger.Error("error : %v", err)

		return "", err
	}

	pocketarkId := data["pocketark_id"].(string)
	ibkId := data["ibk_id"].(string)

	_, err = nk.AccountGetId(ctx, ibkId)

	if err != nil {
		logger.Error("error : %v", err)

		res, _ := json.Marshal(PackMessage{
			"ibk_id": "not_found",
		})

		return string(res), err
	}

	value, err := json.Marshal(PackMessage{
		"pocketark_id": pocketarkId,
		"ibk_id":       ibkId,
	})

	if err != nil {
		logger.Error("error : %v", err)

		return "", err
	}

	_, err = nk.StorageWrite(ctx, []*runtime.StorageWrite{
		{
			Collection: "pocketark",
			Key:        pocketarkId,
			UserID:     ibkId,
			Value:      string(value),
		},
	})

	if err != nil {
		logger.WithField("err", err).Error("StorageWrite error.")
		return "", err
	}

	updated, _, _ := nk.WalletUpdate(ctx, ibkId, ChangSet{
		"pocketark": 1,
	}, PackMessage{"desc": "pocketark"}, true)

	res, _ := json.Marshal(PackMessage{
		"ibk_id": ibkId,
		"count":  updated["pocketark"],
	})

	return string(res), err
}
