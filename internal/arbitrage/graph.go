package arbitrage

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/darth-raijin/challenger/internal/models"
)

type Graph struct {
	Nodes            map[string]map[string]float64
	GoodEnoughProfit float64
}

func NewGraph(goodEnoughProfit float64) *Graph {
	return &Graph{
		Nodes:            make(map[string]map[string]float64),
		GoodEnoughProfit: goodEnoughProfit,
	}
}

func (g *Graph) AddOrder(order models.Order) {
	if g.Nodes[order.BuySymbol] == nil {
		g.Nodes[order.BuySymbol] = make(map[string]float64)
	}
	g.Nodes[order.BuySymbol][order.SellSymbol] = order.Price
}

func (g *Graph) FindArbitrage() ([]string, float64) {
	results := make(chan []string, 100)
	var wg sync.WaitGroup
	var stop int32

	nodes := make([]string, 0, len(g.Nodes))
	for node := range g.Nodes {
		nodes = append(nodes, node)
	}

	batchSize := len(nodes) / 10
	if batchSize == 0 {
		batchSize = 1
	}

	// Create batches of nodes
	for i := 0; i < len(nodes); i += batchSize {
		fmt.Println("Processing batch", i, "to", i+batchSize)
		end := i + batchSize
		if end > len(nodes) {
			end = len(nodes)
		}
		batch := nodes[i:end]
		wg.Add(1)
		go func(batch []string) {
			defer wg.Done()
			for _, start := range batch {
				if atomic.LoadInt32(&stop) == 1 {
					return // Stop processing if good enough profit is found
				}
				distance, previous, hasCycle := bellmanFord(g, start)
				if hasCycle {
					for node, value := range distance {
						if value > 1 {
							cyclePath, _ := traceCycle(node, previous, distance)
							fmt.Println("Found cycle:", cyclePath)
							results <- cyclePath
						}
					}
				}
			}
		}(batch)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	bestProfit := 0.0
	var bestPath []string
	for path := range results {
		if atomic.LoadInt32(&stop) == 1 {
			continue
		}
		profit := calculateProfit(path, g)
		if profit > bestProfit {
			bestProfit = profit
			bestPath = path
			if profit > g.GoodEnoughProfit {
				atomic.StoreInt32(&stop, 1)
				fmt.Println("Terminating early due to high profit:", profit)
			}
		}
	}

	fmt.Println("Best path found:", bestPath)
	return bestPath, bestProfit
}

func calculateProfit(path []string, graph *Graph) float64 {
	if len(path) < 2 {
		return 0
	}

	startCurrency := path[0]
	currentAmount := 1.0

	for i := 0; i < len(path)-1; i++ {
		currentRate := graph.Nodes[path[i]][path[i+1]]
		currentAmount *= currentRate
	}

	finalRate := graph.Nodes[path[len(path)-1]][startCurrency]
	currentAmount *= finalRate

	// Calculate the profit: what we end up with minus what we started with
	return currentAmount - 1 // Started with 1 unit of the start currency
}

func traceCycle(node string, previous map[string]string, distance map[string]float64) ([]string, float64) {
	var path []string
	current := node
	startValue := distance[node]
	visited := make(map[string]bool)

	for {
		if visited[current] {
			break
		}
		visited[current] = true
		path = append([]string{current}, path...)
		current = previous[current]
	}

	if current != node {
		return nil, 0
	}

	profit := startValue - 1
	return path, profit
}

func bellmanFord(g *Graph, start string) (map[string]float64, map[string]string, bool) {
	now := time.Now()
	distance := make(map[string]float64)
	previous := make(map[string]string)

	// Initialize distances
	for node := range g.Nodes {
		distance[node] = 0
	}
	distance[start] = 1 // Assuming '1' is the unit value of the starting node

	// Relax edges up to V-1 times
	for i := 0; i < len(g.Nodes)-1; i++ {
		for u := range g.Nodes {
			for v, weight := range g.Nodes[u] {
				if distance[u] != 0 && distance[u]*weight > distance[v] {
					distance[v] = distance[u] * weight
					previous[v] = u
				}
			}
		}
	}

	// Check for a negative weight cycle
	hasCycle := false
	for u := range g.Nodes {
		for v, weight := range g.Nodes[u] {
			if distance[u] != 0 && distance[u]*weight > distance[v] {
				hasCycle = true
				break
			}
		}
		if hasCycle {
			break
		}
	}

	fmt.Println("Bellman-Ford took", time.Since(now))

	return distance, previous, hasCycle
}
