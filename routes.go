package main

func (s *server) routes() {
	s.router.HandleFunc("/signup", s.handleSignUp)
	s.router.HandleFunc("/signin", s.handleSignIn)
	s.router.HandleFunc("/lists", s.handleLists)
	s.router.HandleFunc("/columns", s.handleColumns)
	s.router.HandleFunc("/items", s.handleItems)
}
