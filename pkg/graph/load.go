package graph

import (
	"encoding/json"
	"os"
)

func NewFromFile(filename string) (*Graph, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var g Graph
	err = json.Unmarshal(bytes, &g)
	if err != nil {
		return nil, err
	}
	return &g, nil
}
