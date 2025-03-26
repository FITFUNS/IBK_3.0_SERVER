package main

import (
	"encoding/json"
)

type PackMessage map[string]interface{}
type ChangSet map[string]int64
type StringDictionary map[string]string
type ObjectMap map[string]json.RawMessage
