module main

import server

fn main() {
	mut s := server.new()!
	for p, ok := s.accept(); ok;{
		p.message("test from v")
	}
}
