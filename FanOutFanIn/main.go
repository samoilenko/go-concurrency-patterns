package main

import (
	"context"
	"fmt"
	"sync"
)

func fanOut(ctx context.Context, inputCh chan int) []chan int {
	numWorkers := 10
	channels := make([]chan int, numWorkers)

	for i := 0; i < numWorkers; i++ {
		addResultCh := add(ctx, inputCh)
		channels[i] = addResultCh
	}

	return channels
}

func fanIn(ctx context.Context, resultChs ...chan int) chan int {
	finalCh := make(chan int)
	var wg sync.WaitGroup

	for _, ch := range resultChs {
		wg.Add(1)
		chClosure := ch

		go func() {
			defer wg.Done()

			for data := range chClosure {
				select {
				case <-ctx.Done():
					return
				case finalCh <- data:
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(finalCh)
	}()

	return finalCh
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

func main() {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	inputCh := generatorWithContext(ctx, input)
	channels := fanOut(ctx, inputCh)
	addResultCh := fanIn(ctx, channels...)

	resultCh := multiply(ctx, addResultCh)

	for res := range resultCh {
		fmt.Println(res)
	}
}
