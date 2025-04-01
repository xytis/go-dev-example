package handlers

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

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
