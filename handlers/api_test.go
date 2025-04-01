package handlers

import (
	"bytes"
	"context"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/xytis/go-dev-example/internal/queue"
)

func TestBufferedPublishConsume(t *testing.T) {
	ctx := context.Background()

	q := queue.NewArrayQueue()
	h := NewHandler(q)

	example := []byte("message")

	req1 := httptest.NewRequestWithContext(ctx, "POST", "/", bytes.NewReader(example))

	wr1 := httptest.NewRecorder()
	h.ServeHTTP(wr1, req1)

	// Note: I specifically did not use any helper libraries.
	//  Usually I just testify/assert or testify/require.
	if wr1.Result().StatusCode != http.StatusOK {
		t.Errorf("unexpected status code: %d", wr1.Result().StatusCode)
	}

	req2 := httptest.NewRequestWithContext(ctx, "GET", "/", nil)

	wr2 := httptest.NewRecorder()
	h.ServeHTTP(wr2, req2)

	if wr2.Result().StatusCode != http.StatusOK {
		t.Errorf("unexpected status code: %d", wr2.Result().StatusCode)
	}

	resp, _ := io.ReadAll(wr2.Result().Body)
	if !bytes.Equal(resp, example) {
		t.Errorf("unexpected response: %s", string(resp))
	}
}

func TestEmptyQueue(t *testing.T) {
	ctx := context.Background()
	q := queue.NewArrayQueue()
	h := NewHandler(q)

	req := httptest.NewRequestWithContext(ctx, "GET", "/", nil)

	wr := httptest.NewRecorder()
	h.ServeHTTP(wr, req)

	if wr.Result().StatusCode != http.StatusNotFound {
		t.Errorf("unexpected status code: %d", wr.Result().StatusCode)
	}
}

func TestInterleaving(t *testing.T) {
	ctx := context.Background()
	q := queue.NewArrayQueue()
	h := NewHandler(q)

	r := rand.New(rand.NewSource(4)) // :)

	var symbols = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	size := 1024 * 1024
	msgLen := 64
	cases := make([][]byte, size)
	for i := 0; i < size; i++ {
		c := make([]byte, msgLen)
		for i := range c {
			c[i] = symbols[byte(r.Intn(len(symbols)))]
		}
		cases[i] = c
	}

	// Start the consumer
	done := make(chan bool)
	go func() {
		defer func() { done <- true }()
		var expected []byte
		current := 0
		for {
			if current >= size {
				break
			}
			expected = cases[current]
			req := httptest.NewRequestWithContext(ctx, "GET", "/", nil)

			wr := httptest.NewRecorder()
			h.ServeHTTP(wr, req)

			if wr.Result().StatusCode == http.StatusNotFound {
				continue
			}
			if wr.Result().StatusCode != http.StatusOK {
				t.Errorf("unexpected status code: %d", wr.Result().StatusCode)
			}

			resp, _ := io.ReadAll(wr.Result().Body)
			if !bytes.Equal(resp, expected) {
				t.Errorf("unexpected response: %s", string(resp))
			}
			current++
		}
	}()

	// Produce
	go func() {
		for i := 0; i < size; i++ {
			req := httptest.NewRequestWithContext(ctx, "POST", "/", bytes.NewReader(cases[i]))

			wr := httptest.NewRecorder()
			h.ServeHTTP(wr, req)

			if wr.Result().StatusCode != http.StatusOK {
				t.Errorf("unexpected status code: %d", wr.Result().StatusCode)
			}
		}
	}()

	// Block
	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Errorf("timeout")
	}
}
