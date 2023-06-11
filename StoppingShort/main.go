package main

import (
	"context"
	"fmt"
	"time"
)

func handler() {
	input := []int{1, 2, 3, 4, 5, 6}

	// Explicit cancellation
	doneCh := make(chan struct{})
	defer close(doneCh)

	inputCh := generator(doneCh, input)

	for data := range inputCh {
		if data == 1 {
			return
		}
	}
}

func generator(doneCh chan struct{}, input []int) chan int {
	inputCh := make(chan int)

	go func() {
		defer close(inputCh)

		for _, data := range input {
			select {
			case <-doneCh:
				return
			case inputCh <- data:
			}
		}
	}()

	return inputCh
}

func handlerWithContext(ctx context.Context) {
	input := []int{1, 2, 3, 4, 5, 6}
	inputCh := generatorWithContext(ctx, input)

	for data := range inputCh {
		fmt.Println(data)
		if data == 1 {
			return
		}
	}
}

func generatorWithContext(ctx context.Context, input []int) chan int {
	inputCh := make(chan int)

	go func() {
		defer close(inputCh)

		for _, data := range input {
			select {
			case <-ctx.Done():
				fmt.Println("Exit generator")
				return
			case inputCh <- data:
			}
		}
	}()

	return inputCh
}

func main() {
	handler()

	ctx, cancel := context.WithCancel(context.Background())
	handlerWithContext(ctx)

	// Explicit cancellation

	defer func() {
		cancel()
		time.Sleep(5 * 500 * time.Millisecond)
	}()
}
