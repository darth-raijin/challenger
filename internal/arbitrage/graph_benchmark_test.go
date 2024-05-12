package arbitrage_test

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/darth-raijin/challenger/internal/arbitrage"
	"github.com/darth-raijin/challenger/internal/models"
	"github.com/google/uuid"
)

func BenchmarkGraph_FindArbitrage(b *testing.B) {
	for i := 1; i <= runtime.NumCPU(); i++ {
		runtime.GOMAXPROCS(i)
	}

	for i := 0; i < 1; i++ {
		fmt.Println("iteration: ", i, "of", b.N, "benchmarks")
		graph := arbitrage.NewGraph()

		for j := 0; j < 500; j++ {
			tickerOne := uuid.NewString()
			tickerTwo := uuid.NewString()
			tickerThree := uuid.NewString()
			tickerFour := uuid.NewString()

			graph.AddOrder(models.Order{BuySymbol: tickerOne, SellSymbol: tickerTwo, Price: 50})
			graph.AddOrder(models.Order{BuySymbol: tickerOne, SellSymbol: tickerThree, Price: 25})
			graph.AddOrder(models.Order{BuySymbol: tickerTwo, SellSymbol: tickerFour, Price: 2000})
			graph.AddOrder(models.Order{BuySymbol: tickerThree, SellSymbol: tickerFour, Price: 500})
			graph.AddOrder(models.Order{BuySymbol: tickerFour, SellSymbol: tickerOne, Price: 0.0011})
		}

		b.ResetTimer()
		fmt.Println("Finding arbitrage")
		graph.FindArbitrage()
	}
}
