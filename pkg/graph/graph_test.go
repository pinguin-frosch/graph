package graph_test

import (
	"graph/pkg/graph"
	"testing"
)

func TestNewNode(t *testing.T) {
	invalidIds := []string{"normal1", "thisisfine()", "...777", "_._asdfñ"}
	for _, invalidId := range invalidIds {
		_, err := graph.NewNode(invalidId)
		if err == nil {
			t.Fatalf(`NewNode("%v") should return an error, %s is invalid`, invalidId, invalidId)
		}
	}
	validIds := []string{"ABC__jj", "jklIIO.aa._z", "a_b_c_d", "a.z.i.n"}
	for _, validId := range validIds {
		_, err := graph.NewNode(validId)
		if err != nil {
			t.Fatalf(`NewNode("%v") should'nt return an error, %s is valid`, validId, validId)
		}
	}
}
