package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/xytis/go-dev-example/handlers"
	"github.com/xytis/go-dev-example/internal/queue"
)

func main() {
	// Note: I intentionally did not use cobra or similar, and kept the initialization
	//  in the main function, instead of hoisting into init().

	// For actual implementation I would use library like cobra, and split different operating
	//  modes under different subcommands ('_ consume', '_ produce' and '_ queue')

	consumer := flag.Bool("queue", false, "Queue mode")

	input := flag.String("input", "in.txt", "Input file")
	output := flag.String("output", "out.txt", "Output file")

	flag.Parse()

	if *consumer {
		runQueue()
		return
	}

	runClient(*input, *output)
}

func runQueue() {
	fmt.Println("Queue mode")

	q := queue.NewArrayQueue()

	h := handlers.NewHandler(q)

	if err := http.ListenAndServe(":8080", h); err != nil {
		// Note: panic is intentional here.
		//  Even though in this particular case it will provide zero extra information
		//  during runtime, I used it here to denote that in production extra code should
		//  wrap a http server. From the top of my head -> graceful shutdown should be implemented.
		panic(err)
	}
}

func runClient(input, output string) {
	fmt.Println(input)
	fmt.Println(output)
}
