package main

import (
	"context"
	"database/sql"

	"github.com/heroiclabs/nakama-common/runtime"
)

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	logger.Info("Hello World!")

	nk.LeaderboardCreate(ctx, AUGUST_LEADERBOARD, true, "desc", "set", "", map[string]interface{}{}, true)
	nk.LeaderboardCreate(ctx, SEPTEMBER_LEADERBOARD, true, "desc", "set", "", map[string]interface{}{}, true)
	nk.LeaderboardCreate(ctx, OCTOBER_LEADERBOARD, true, "desc", "set", "", map[string]interface{}{}, true)
	nk.LeaderboardCreate(ctx, GENERAL_LEADERBOARD, true, "desc", "set", "", map[string]interface{}{}, true)

	initializer.RegisterRpc("remove", RemoveAccount)
	initializer.RegisterRpc("get_rank", GetRank)
	initializer.RegisterRpc("submit_water_rank", SubmitWaterRank)
	initializer.RegisterRpc("submit_gummy_rank", SubmitGummyRank)
	initializer.RegisterRpc("submit_melon_rank", SubmitMelonRank)
	initializer.RegisterRpc("get_wallet", GetWallet)
	initializer.RegisterRpc("get_fitfuns_reward", GetFitfunsReward)
	initializer.RegisterRpc("update_privacy", UpdatePrivacy)
	initializer.RegisterRpc("update_allow", UpdateAllow)
	initializer.RegisterRpc("get_privacy", GetPrivacy)
	initializer.RegisterRpc("get_calinder", GetCalinder)
	initializer.RegisterRpc("get_pocketark_id", GetPocketarkId)

	return nil
}
