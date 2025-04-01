package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

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

	fmt.Println(" listening on :8080")
	if err := http.ListenAndServe(":8080", h); err != nil {
		// Note: panic is intentional here.
		//  Even though in this particular case it will provide zero extra information
		//  during runtime, I used it here to denote that in production extra code should
		//  wrap a http server. From the top of my head -> graceful shutdown should be implemented.
		panic(err)
	}
}

func runClient(input, output string) {
	ctx := context.Background()

	{
		fmt.Println("reading from ", input, " writing to http://localhost:8080")
		in, err := os.Open(input)
		if err != nil {
			panic(err)
		}
		defer in.Close()
		// Note: defers are not executed on lexical scope end. This is purely for aesthetics.

		scanner := bufio.NewScanner(in)
		for scanner.Scan() {
			msg := scanner.Text()
			req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:8080/", strings.NewReader(msg))
			if err != nil {
				panic(err)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				panic(err)
			}
			resp.Body.Close()
		}

		if err := scanner.Err(); err != nil {
			panic(err)
		}

		fmt.Println("done")
	}

	{
		fmt.Println("reading from http://localhost:8080 writing to ", output)
		out, err := os.Create(output)
		if err != nil {
			panic(err)
		}
		defer out.Close()

		for {
			req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/", nil)
			if err != nil {
				panic(err)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
				panic("unexpected status code")
			}
			if resp.StatusCode == http.StatusNotFound {
				break
			}

			msg, err := io.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}
			_, err = fmt.Fprintln(out, string(msg))
			if err != nil {
				panic(err)
			}
		}

		fmt.Println("done")
	}
}
