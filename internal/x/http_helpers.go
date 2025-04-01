package x

import (
	"io"
	"net/http"
)

// ApiError is a stub, which should help with standardization of API responses.
//
//	In real application a framework would be used instead.
func ApiError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func ApiSuccess(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func ApiEmpty(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

func ApiPayload(w http.ResponseWriter, payload []byte) {
	w.WriteHeader(http.StatusOK)
	n, err := w.Write(payload)
	// Note: the below checks are not too useful in this demo example,
	//  also this behaviour is out of scope for this demo task.
	//  (that is transport framework domain)
	if err != nil {
		ApiError(w, err)
		return
	}
	if n != len(payload) {
		ApiError(w, io.ErrShortWrite)
		return
	}
}
