package graph

import (
	"encoding/json"
	"errors"
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

func (g *Graph) saveSnapshot() error {
	bytes, err := json.Marshal(g)
	if err != nil {
		return err
	}
	err = os.Mkdir("snapshots", 0755)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}
	filename := fmt.Sprintf("snapshots/snapshot.json")
	err = os.WriteFile(filename, bytes, 0644)
	if err != nil {
		return err
	}
	fmt.Printf("Saved as %s\n", filename)
	return nil
}
