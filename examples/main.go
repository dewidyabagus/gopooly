package main

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"time"

	gopooly "github.com/dewidyabagus/gopooly"
)

const counter int = 1234

type Balok struct {
	Panjang int
	Lebar   int
}

func main() {
	// runtime.GOMAXPROCS(2)

	var execFunc = func(ctx context.Context, args interface{}) (interface{}, error) {
		values, ok := args.(Balok)
		if !ok {
			return nil, errors.New("data type not Balok")
		}
		luas := values.Panjang * values.Lebar
		fmt.Println("Luas Balok: P x L =", values.Panjang, "*", values.Lebar, "=", luas)
		time.Sleep(time.Millisecond * 50)
		return values, nil
	}

	t := gopooly.NewFunc(2000, 10000, execFunc)
	defer func() {
		time.Sleep(time.Millisecond)

		fmt.Println("Total Goroutine:", runtime.NumGoroutine())
	}()

	start := time.Now()

	// wg := new(sync.WaitGroup)

	for i := 0; i < counter; i++ {
		t.Process(context.Background(), Balok{Panjang: i + 5, Lebar: i})
		// t.ExecFunc(context.Background(), Balok{Panjang: i + 5, Lebar: i})

		// Meluncurkan langsung
		// wg.Add(1)
		// go func(n int) {
		// 	defer wg.Done()

		// 	execFunc(context.Background(), Balok{Panjang: n + 5, Lebar: n})
		// }(i)
	}

	// wg.Wait()
	t.Close()
	stop := time.Now()

	time.Sleep(time.Second)
	fmt.Println("Start At   :", start)
	fmt.Println("Latency    :", stop.Sub(start))
	fmt.Println("Len queue  :", t.Len())
	fmt.Println("Cap queue  :", t.Cap())
	fmt.Println("Task Total :", t.TaskTotal())
}
