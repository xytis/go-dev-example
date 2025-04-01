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

func port() string {
	p, found := os.LookupEnv("PORT")
	if !found {
		p = "8080"
	}
	return p
}

func main() {
	// Note: I intentionally did not use cobra or similar, and kept the initialization
	//  in the main function, instead of hoisting into init().

	// For actual implementation I would use library like cobra, and split different operating
	//  modes under different subcommands ('_ consume', '_ produce' and '_ queue')

	consumer := flag.Bool("queue", false, "Queue mode")

	input := flag.String("input", "in.txt", "Input file")
	output := flag.String("output", "out.txt", "Output file")
	url := flag.String("url", "http://localhost:"+port(), "Remote queue URL")

	flag.Parse()

	if *consumer {
		runQueue()
		return
	}

	runClient(*url, *input, *output)
}

func runQueue() {
	fmt.Println("Queue mode")

	q := queue.NewArrayQueue()

	h := handlers.NewHandler(q)

	fmt.Printf("  listening on :%s\n", port())
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port()), h); err != nil {
		// Note: panic is intentional here.
		//  Even though in this particular case it will provide zero extra information
		//  during runtime, I used it here to denote that in production extra code should
		//  wrap a http server. From the top of my head -> graceful shutdown should be implemented.
		panic(err)
	}
}

func runClient(url, input, output string) {
	ctx := context.Background()

	clientUpload(ctx, url, input)
	clientDownload(ctx, url, output)
}

func clientUpload(ctx context.Context, url, input string) {
	fmt.Println("reading from", input, "writing to", url)
	in, err := os.Open(input)
	if err != nil {
		panic(err)
	}
	defer in.Close()
	// Note: defers are not executed on lexical scope end. This is purely for aesthetics.

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		msg := scanner.Text()
		req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(msg))
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

func clientDownload(ctx context.Context, url, output string) {
	fmt.Println("reading from", url, "writing to", output)
	out, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	fetchOne := func() ([]byte, bool, error) {
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, false, err
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, false, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
			return nil, false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}
		if resp.StatusCode == http.StatusNotFound {
			return nil, true, nil
		}

		msg, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, false, err
		}
		return msg, false, err
	}

	for {
		msg, done, err := fetchOne()
		if err != nil {
			panic(err)
		}
		if done {
			break
		}

		_, err = fmt.Fprintln(out, string(msg))
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("done")
}
