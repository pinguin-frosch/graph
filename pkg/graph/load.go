package graph

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

type GraphFile []string

func LoadFromJson(filename string) (*Graph, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var file GraphFile
	err = json.Unmarshal(bytes, &file)
	if err != nil {
		return nil, err
	}
	g := Graph{}
	for _, e := range file {
		parts := strings.Split(e, "|")
		if len(parts) != 2 {
			return nil, errors.New(fmt.Sprintf("Invalid edge: %v", e))
		}
		from := parts[0]
		err := g.AddVertex(from)
		if err != nil && err.Error() != RepeatedVertex {
			return nil, err
		}
		tos := strings.Split(parts[1], ".")
		for _, to := range tos {
			deadEnd := false
			if strings.HasSuffix(to, "*") {
				to = strings.TrimSuffix(to, "*")
				deadEnd = true
			}
			err := g.AddVertex(to)
			if err != nil && err.Error() != RepeatedVertex {
				return nil, err
			}
			err = g.AddEdge(from, to)
			if err != nil && err.Error() != RepeatedEdge {
				return nil, err
			}
			if deadEnd {
				f := g.GetEdge(from, to)
				if f == nil {
					return nil, errors.New(fmt.Sprintf("Couldn't get the edge %v<->%v", from, to))
				}
				f.deadEnd = true
			}
		}
	}
	return &g, nil
}
