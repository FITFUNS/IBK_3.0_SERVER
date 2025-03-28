package main

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
)

type MatchState struct {
	presences map[string]Player
	messages  []interface{}
}

type Match struct{}

func newMatch(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule) (m runtime.Match, err error) {
	return &Match{}, nil
}

func (m *Match) MatchInit(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, params map[string]interface{}) (interface{}, int, string) {
	state := &MatchState{
		presences: make(map[string]Player),
		messages:  make([]interface{}, 0),
	}

	tickRate := TICK_RATE
	label := MATCH_LABEL

	return state, tickRate, label
}

func (m *Match) MatchJoinAttempt(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presence runtime.Presence, metadata map[string]string) (interface{}, bool, string) {
	mState, _ := state.(*MatchState)
	acceptUser := true
	if _, ok := mState.presences[presence.GetUsername()]; ok {
		acceptUser = false
	}
	return state, acceptUser, ""
}

func (m *Match) MatchJoin(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presences []runtime.Presence) interface{} {
	mState, _ := state.(*MatchState)
	payload := []interface{}{}

	for _, p := range presences {
		account, _ := nk.AccountGetId(ctx, p.GetUserId())
		data := PackMessage{}
		err := json.Unmarshal([]byte(account.User.Metadata), &data)

		if err != nil {
			continue
		}

		player := newPlayer(p, account.User.DisplayName, data["pos"])

		mState.presences[p.GetUserId()] = *player
		payload = append(payload, player.GetPayload())
	}
	initPayload := make([]interface{}, 0)
	for _, p := range mState.presences {
		initPayload = append(initPayload, p.GetPayload())
	}
	jsonData, err := json.Marshal(initPayload)
	if err != nil {
		logger.Error("Error marshalling messages: %v", err)
		return mState
	}
	dispatcher.BroadcastMessageDeferred(OP_CODE_INIT, jsonData, presences, nil, false)
	mState.messages = append(mState.messages, []interface{}{OP_CODE_SPAWN, payload})
	return mState
}

func (m *Match) MatchLeave(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presences []runtime.Presence) interface{} {
	mState, _ := state.(*MatchState)

	for _, p := range presences {
		player, ok := mState.presences[p.GetUserId()]

		if !ok {
			break
		}
		account, _ := nk.AccountGetId(ctx, p.GetUserId())
		data := map[string]interface{}{}
		err := json.Unmarshal([]byte(account.User.Metadata), &data)

		if err != nil {
			continue
		}

		data["pos"] = player.Position

		nk.AccountUpdateId(ctx, p.GetUserId(), "", data, "", "", "", "", "")
		delete(mState.presences, p.GetUserId())
	}

	return mState
}

func (m *Match) MatchLoop(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, messages []runtime.MatchData) interface{} {
	mState, _ := state.(*MatchState)

	for _, message := range messages {
		HandleMessage(message, mState)
	}

	if len(mState.messages) > 0 {
		jsonData, err := json.Marshal(mState.messages)
		if err != nil {
			logger.Error("Error marshalling messages: %v", err)
			return mState
		}
		// logger.Info("Broadcasting %v messages", len(mState.messages))
		dispatcher.BroadcastMessageDeferred(OP_CODE_PACK, jsonData, nil, nil, false)
		mState.messages = mState.messages[:0]
	}

	return mState
}

func (m *Match) MatchTerminate(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, graceSeconds int) interface{} {
	return state
}

func (m *Match) MatchSignal(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, data string) (interface{}, string) {
	return state, "signal received: " + data
}
