package webscan

import (
	"embed"
	"encoding/json"
)

//go:embed data/fingers.json
var f embed.FS

func GetDefaultFingersData() (fingerRules []*FingerRule, err error) {
	bytes, err := f.ReadFile("data/fingers.json")
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, &fingerRules)
	return
}
