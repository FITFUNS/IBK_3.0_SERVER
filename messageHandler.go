package main

import (
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
)

func HandleMessage(message runtime.MatchData, mState *MatchState) {
	opCode := message.GetOpCode()
	data := message.GetData()
	messageSender := message.GetUsername()
	switch opCode {
	case OP_CODE_MOVE:
		HandleMove(data, mState, messageSender)
	}
}
func HandleMove(data []byte, mState *MatchState, messageSender string) {
	payload := MovePayload{}
	err := json.Unmarshal(data, &payload)
	if err != nil {
		return
	}

	player, ok := mState.presences[messageSender]
	if !ok {
		return
	}

	player.Position = payload

	mState.messages = append(mState.messages, player.GetPayload())
}
