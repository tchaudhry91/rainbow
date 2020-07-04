package service

// routes registers the handlers to the specific routes
func (s *RainbowHTTP) routes() {
	s.router.HandleFunc("/hash", s.GetHash()).Methods("GET")
	s.router.HandleFunc("/reverse", s.GetReverseHash()).Methods("GET")
}
