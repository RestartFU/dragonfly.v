module main

import server

fn main() {
	mut s := server.new()!
	for {
		mut p, ok := s.accept()
		if !ok {
			return
		}
		p.message("welcome, ${p.name()}")
	}
}
