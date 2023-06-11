package main

import (
	"context"
	"fmt"
)

func add(ctx context.Context, inputCh chan int) chan int {
	addRes := make(chan int)

	go func() {
		defer func() {
			fmt.Println("Exit add")
			close(addRes)
		}()

		for data := range inputCh {
			result := data + 1

			select {
			case <-ctx.Done():
				fmt.Println("Exit add ctx")
				return
			case addRes <- result:
			}
		}
	}()

	return addRes
}

func multiply(ctx context.Context, inputCh chan int) chan int {
	multiplyRes := make(chan int)

	go func() {
		defer func() {
			fmt.Println("Exit multiply")
			close(multiplyRes)
		}()

		for data := range inputCh {
			result := data * 2

			select {
			case <-ctx.Done():
				fmt.Println("Exit multiply ctx")
				return
			case multiplyRes <- result:
			}
		}
	}()

	return multiplyRes
}

func generatorWithContext(ctx context.Context, input []int) chan int {
	inputCh := make(chan int)

	go func() {
		defer func() {
			fmt.Println("Exit generator")
			close(inputCh)
		}()

		for _, data := range input {
			select {
			case <-ctx.Done():
				fmt.Println("Exit generator context")
				return
			case inputCh <- data:
			}
		}
	}()

	return inputCh
}

func main() {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8}
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		// time.Sleep(5 * 100 * time.Millisecond)
	}()

	inputCh := generatorWithContext(ctx, input)
	resultCh := multiply(ctx, add(ctx, inputCh))

	for res := range resultCh {
		// using cancel() in this app is useless, because all channels will be closed after numbers in generator is ended
		// but the simulation is to run cancel() below
		// cancel()
		fmt.Println(res)
	}

}
