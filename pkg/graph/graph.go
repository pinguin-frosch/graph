package graph

type Graph struct {
	vertices []*Vertex
}

func (g *Graph) AddVertex(key string) {
	v := &Vertex{
		key: key,
	}
	g.vertices = append(g.vertices, v)
}

type Vertex struct {
	key string
}
