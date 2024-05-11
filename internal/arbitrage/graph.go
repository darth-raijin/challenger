package arbitrage

import (
	"fmt"
	"math"

	"github.com/darth-raijin/challenger/internal/models"
)

type Graph struct {
	Nodes map[string]map[string]float64
}

func NewGraph() *Graph {
	return &Graph{Nodes: make(map[string]map[string]float64)}
}

func (g *Graph) AddOrder(order models.Order) {
	if g.Nodes[order.BuySymbol] == nil {
		g.Nodes[order.BuySymbol] = make(map[string]float64)
	}
	g.Nodes[order.BuySymbol][order.SellSymbol] = order.Price
}

func (g *Graph) FindArbitrage() (path []string, profit float64) {
	for start := range g.Nodes {
		bellmanFord(g, start)

	}
	return nil, 0
}

func bellmanFord(g *Graph, start string) (map[string]float64, map[string]string, bool) {
	distance := make(map[string]float64)
	previous := make(map[string]string)
	for node := range g.Nodes {
		distance[node] = math.Inf(1)
	}
	distance[start] = 1

	for current, next := range g.Nodes[start] {
		fmt.Println(fmt.Sprintf("%v: %v", current, next))
	}

	return distance, previous, false
}
