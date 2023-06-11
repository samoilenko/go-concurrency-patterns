package main

import "log"

type Result struct {
	data int
	err  error
}

func main() {
	input := []int{1, 2, 3, 4}
	resultCh := make(chan Result)

	inputCh := generator(input)

	go consumer(inputCh, resultCh)

	for res := range resultCh {
		if res.err != nil {
			log.Println("deal with the error here")
		}
	}
}

func consumer(inputCh chan int, resultCh chan Result) {
	defer close(resultCh)

	for data := range inputCh {
		resp, err := callDatabase(data)

		result := Result{
			data: resp,
			err:  err,
		}

		resultCh <- result
	}
}
