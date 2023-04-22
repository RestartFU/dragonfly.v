module main

import server
import player

struct Handler {
mut:
	p &player.Player
}

fn (mut h Handler) handle_quit() {
	println('player ${h.p.name()} quit')
}

fn main() {
	mut s := server.new()!
	for {
		mut p, ok := s.accept()
		if !ok {
			return
		}
		p.handle(&Handler{ p: p })
		p.message('welcome, ${p.name()}')
		mut player_list := []string{}
		for mut pl in s.players() {
			player_list << pl.name()
		}
		println('players online: ${player_list.join(', ')}')
	}
}
