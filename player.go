package main

import (
	"math"
	"math/rand"

	"github.com/heroiclabs/nakama-common/runtime"
)

// type Mover interface {

// }

type Player struct {
	UserId      string    `json:"user_id"`
	SessionId   string    `json:"session_id"`
	DisplayName string    `json:"display_name"`
	Username    string    `json:"username"`
	Position    []float64 `json:"position"`
	Direction   []float64 `json:"direction"`
	Speed       float64   `json:"speed"`
	Before      []float64 `json:"before"`
	Node        string    `json:"node"`
	Stopped     bool      `json:"stopped"`
}

func newPlayer(presence runtime.Presence, displayName string, pos interface{}) *Player {
	if pos == nil {
		pos = []float64{
			float64(rand.Intn(3) + 13),
			float64(rand.Intn(3) - 1),
		}
	}
	player := &Player{
		UserId:      presence.GetUserId(),
		SessionId:   presence.GetSessionId(),
		DisplayName: displayName,
		Username:    presence.GetUsername(),
		Position:    pos.([]float64),
		Direction:   []float64{0, 0},
		Speed:       4,
		Before:      pos.([]float64),
		Node:        presence.GetNodeId(),
		Stopped:     false,
	}

	return player
}

func (p *Player) GetPayload() interface{} {
	return []interface{}{
		p.UserId,
		p.DisplayName,
		p.Position,
		p.Direction,
	}
}

func (p *Player) CheckCollision() interface{} {
	isStuck := false

	for _, circle := range MAP_INFO {
		x := circle.Position[0] - p.Position[0]
		z := circle.Position[1] - p.Position[1]

		dist := math.Sqrt(x*x + z*z)

		if dist <= circle.Radius+0.5 {
			p.Position = p.Before
			p.Stopped = false
			isStuck = true
			break
		}
	}

	if isStuck {
		return p.GetPayload()
	}

	if p.Stopped {
		p.Stopped = false
		return p.GetPayload()
	}

	return nil
}

func (p *Player) Update(diff float64) interface{} {
	p.Before = p.Position

	mag := math.Sqrt(p.Direction[0]*p.Direction[0] + p.Direction[1]*p.Direction[1])

	if mag > 0 {
		tempVec := []float64{
			p.Direction[0] / mag,
			p.Direction[1] / mag,
		}
		speed := p.Speed
		p.Position[0] += tempVec[0] * diff * speed
		p.Position[1] += tempVec[1] * diff * speed
		posMag := math.Sqrt(p.Position[0]*p.Position[0] + p.Position[1]*p.Position[1])

		if posMag > MAX_DISTANCE {
			p.Position[0] /= posMag * MAX_DISTANCE
			p.Position[1] /= posMag * MAX_DISTANCE
		}

		return p.CheckCollision()
	}

	if p.Stopped {
		p.Stopped = false
		return p.GetPayload()
	}

	return nil
}
