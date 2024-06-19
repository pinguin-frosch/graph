package graph

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
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
	now := time.Now().Format("02-01-06 15:04:05")
	filename := fmt.Sprintf("snapshots/%s.json", now)
	err = os.WriteFile(filename, bytes, 0644)
	if err != nil {
		return err
	}
	fmt.Printf("Saved as %s\n", filename)
	return nil
}

func (g *Graph) ModifyInteractively() error {
	fmt.Println("? to list the options")

	stdin := bufio.NewReader(os.Stdin)
	exit := false
	for {
		fmt.Printf(">> ")

		var option string
		fmt.Scanf("%s", &option)

		switch option {
		case "A":
			// TODO: add more algorithms
			d := DefaultTraverseAlgorithm{}
			g.SetTraverseAlgorithm(d)
		case "E":
			var from, to, deadEnd string
			var weight uint

			fmt.Printf("From: ")
			fmt.Scanf("%s", &from)

			fmt.Printf("To: ")
			fmt.Scanf("%s", &to)

			fmt.Printf("Deadend [N/y]: ")
			fmt.Scanf("%s", &deadEnd)

			fmt.Printf("Weight (default 0): ")
			n, _ := fmt.Scanf("%d", &weight)
			if n == 0 {
				fmt.Println("Using default value for weight.")
				stdin.ReadString('\n')
			}

			err := g.AddEdge(from, to)
			if err != nil {
				printErrorMessage(err)
				break
			}

			edge := g.GetEdge(from, to)
			if edge == nil {
				fmt.Printf("Error: Couldn't get edge!")
				exit = true
				break
			}
			otherEdge := g.GetEdge(to, from)
			if otherEdge == nil {
				fmt.Printf("Error: Couldn't the other get edge!")
				break
			}

			edge.Weight = weight
			otherEdge.Weight = weight
			if deadEnd == "y" {
				edge.DeadEnd = true
			}

		case "F":
			sequence, err := g.GetShortestSequence()
			if err != nil {
				printErrorMessage(err)
				break
			}
			sequence.Print()

		case "L":
			fmt.Println("After loading the state, any unsaved changes will be lost.")
			var filename string
			fmt.Printf("Filename: ")
			fmt.Scanf("%s", &filename)

			newGraph, err := NewFromFile(filename)
			if err != nil {
				printErrorMessage(err)
				break
			}

			g = newGraph

		case "P":
			g.Print()

		case "S":
			err := g.saveSnapshot()
			if err != nil {
				fmt.Printf("Error: %v\n", err.Error())
			}

		case "V":
			fmt.Println("Adding vertices, type .exit to go back")
			count := 0
			for {
				var vertex string
				fmt.Printf("Vertex: ")
				fmt.Scanf("%s", &vertex)
				if vertex == ".exit" {
					break
				}
				err := g.AddVertex(vertex)
				if err != nil {
					printErrorMessage(err)
					continue
				}
				count++
			}
			if count > 0 {
				fmt.Printf("Succesfully added %v vertices.\n", count)
			}

		case "W":
			var from string
			fmt.Printf("From: ")
			fmt.Scanf("%s", &from)

			sequence, err := g.GetSequence(from)
			if err != nil {
				printErrorMessage(err)
				break
			}
			sequence.Print()

		case "X":
			exit = true

		case "?":
			printOptions()
		}

		if exit {
			break
		}
	}

	return nil

}

func printOptions() {
	fmt.Println("A: Change traverse algorithm")
	fmt.Println("E: Add edge")
	fmt.Println("F: Find shortest path")
	fmt.Println("L: Load graph state from file")
	fmt.Println("P: Print graph")
	fmt.Println("S: Save graph state")
	fmt.Println("V: Add vertices")
	fmt.Println("W: Walk the graph")
	fmt.Println("X: Exit")
	fmt.Println("?: List options")
}

func printErrorMessage(err error) {
	fmt.Printf("Error: %v\n", err.Error())
}
