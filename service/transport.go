package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

// RainbowHTTP is an HTTP server for the underlying RainbowService
type RainbowHTTP struct {
	svc    RainbowService
	server *http.Server
	router *mux.Router
}

// NewRainbowHTTP returns a new HTTP server for the Rainbow Service.
func NewRainbowHTTP(listenAddr string, svc RainbowService) *RainbowHTTP {
	router := mux.NewRouter()

	// Build the server
	server := &RainbowHTTP{
		svc:    svc,
		router: router,
		server: &http.Server{
			Addr:    listenAddr,
			Handler: router,
		},
	}

	// Initialize the routes
	server.routes()

	return server
}

// Handler returns the underlying HTTP handler
func (s *RainbowHTTP) Handler() http.Handler {
	return s.router
}

// Start begins listening for requests on the bindAddr. Blocks.
func (s *RainbowHTTP) Start() error {
	return s.server.ListenAndServe()
}

// Shutdown gracefully terminates the server
func (s *RainbowHTTP) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// respond is a internal utility to set proper HTTP responses
func (s *RainbowHTTP) respond(w http.ResponseWriter, req *http.Request, data interface{}, statusCode int, err error) {
	w.WriteHeader(statusCode)
	// Log the supplied error later if it's not nil
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			// Log this later
		}
	}
}

// GetHash returns the HandlerFunc for the Hash route.
func (s *RainbowHTTP) GetHash() http.HandlerFunc {
	// Do some potential setup here.
	// Not required in this case.

	return func(w http.ResponseWriter, req *http.Request) {
		// These structs are limited to the handler scope.
		// They may seem overly verbose in this case, but they're worth it when you're dealing with more complex requests.
		type Request struct {
			Str string `json:"str,omitempty"`
		}
		type Response struct {
			Hash string `json:"hash,omitempty"`
			Err  string `json:"err,omitempty"`
		}

		r := Request{}

		qparams := req.URL.Query()
		if _, ok := qparams["str"]; !ok {
			s.respond(w, req, nil, http.StatusBadRequest, errors.New("No string supplied"))
			return
		}

		// Grab the first "str" query parameter
		r.Str = qparams["str"][0]

		resp := Response{}

		hash, err := s.svc.Hash(r.Str)
		if err != nil {
			resp.Err = err.Error()
			s.respond(w, req, resp, http.StatusInternalServerError, err)
			return
		}
		resp.Hash = hash
		s.respond(w, req, resp, http.StatusOK, nil)

	}
}

// GetReverseHash returns the handlerFunc for the Get ReverseHash route
func (s *RainbowHTTP) GetReverseHash() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		type Request struct {
			Hash string `json:"str,omitempty"`
		}

		type Response struct {
			Str string `json:"str,omitempty"`
			Err string `json:"err,omitempty"`
		}

		r := Request{}

		qparams := req.URL.Query()
		if _, ok := qparams["hash"]; !ok {
			s.respond(w, req, nil, http.StatusBadRequest, errors.New("No hash supplied"))
			return
		}

		// Grab the first "hash" query parameter
		r.Hash = qparams["hash"][0]

		resp := Response{}

		hash, err := s.svc.Hash(r.Hash)
		if err != nil {
			resp.Err = err.Error()
			s.respond(w, req, resp, http.StatusInternalServerError, err)
			return
		}
		resp.Str = hash
		s.respond(w, req, resp, http.StatusOK, nil)
	}
}
