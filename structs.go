package main

import (
	"encoding/json"
)

type PackMessage map[string]interface{}
type ChangSet map[string]int64
type StringDictionary map[string]string
type ObjectMap map[string]json.RawMessage
type MovePayload []float64

type Vector2 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
