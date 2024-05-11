package arbitrage_test

import (
	"testing"

	"github.com/darth-raijin/challenger/internal/arbitrage"
	"github.com/darth-raijin/challenger/internal/models"
)

// TestNewGraph tests the creation of a new Graph
func TestNewGraph(t *testing.T) {
	graph := arbitrage.NewGraph()
	if graph == nil {
		t.Error("NewGraph() = nil, expected non-nil graph instance")
	}
	if len(graph.Nodes) != 0 {
		t.Errorf("New graph should have no nodes, got %v", len(graph.Nodes))
	}
}

// TestAddOrder tests adding orders to the graph
func TestAddOrder(t *testing.T) {
	graph := arbitrage.NewGraph()
	order := models.Order{
		BuySymbol:  "BTC",
		SellSymbol: "ETH",
		Price:      50.0, // Example price
	}
	graph.AddOrder(order)

	if len(graph.Nodes["BTC"]) == 0 {
		t.Error("AddOrder did not add the order correctly")
	}
	if graph.Nodes["BTC"]["ETH"] != 50.0 {
		t.Errorf("Expected BTC to ETH rate to be 50.0, got %v", graph.Nodes["BTC"]["ETH"])
	}
}

// TestFindArbitrage with no arbitrage scenario
func TestFindArbitrage_NoArbitrage(t *testing.T) {
	graph := arbitrage.NewGraph()
	graph.AddOrder(models.Order{BuySymbol: "BTC", SellSymbol: "ETH", Price: 10})
	graph.AddOrder(models.Order{BuySymbol: "ETH", SellSymbol: "BTC", Price: 0.1})

	_, profit := graph.FindArbitrage()
	if profit != 0 {
		t.Errorf("Expected no arbitrage, got profit %v", profit)
	}
}

// TestFindArbitrage with an arbitrage scenario
func TestFindArbitrage_WithArbitrage(t *testing.T) {
	graph := arbitrage.NewGraph()
	graph.AddOrder(models.Order{BuySymbol: "BTC", SellSymbol: "ETH", Price: 50})
	graph.AddOrder(models.Order{BuySymbol: "BTC", SellSymbol: "SOL", Price: 25})
	graph.AddOrder(models.Order{BuySymbol: "ETH", SellSymbol: "USD", Price: 2000})
	graph.AddOrder(models.Order{BuySymbol: "SOL", SellSymbol: "USD", Price: 500})
	graph.AddOrder(models.Order{BuySymbol: "USD", SellSymbol: "BTC", Price: 0.0011}) // Creates arbitrage

	_, profit := graph.FindArbitrage()
	if profit <= 1 {
		t.Errorf("Expected profitable arbitrage, got profit %v", profit)
	}
}
