module server

import dl.loader
import player
import arrays

type T_server_start = fn() voidptr
type T_server_accept = fn(voidptr) voidptr
type T_server_players = fn(voidptr) voidptr

pub struct Server{
	handle voidptr

	server_accept T_server_accept
	server_players T_server_players

	mut: players map[voidptr]&player.Player
}

pub fn new() !&Server{
	mut libs := []string{}
	libs << "dragonfly.dll"

	mut lib := loader.get_or_create_dynamic_lib_loader(loader.DynamicLibLoaderConfig{paths: libs}) or {
		panic(err)
	}

	lib.open()!

	start_server := T_server_start(lib.get_sym("server_start")!)
	h := start_server()

	mut s := &Server{
		handle: h

		server_accept: T_server_accept(lib.get_sym("server_accept")!)
		server_players: T_server_players(lib.get_sym("server_players")!)
	}

	return s
}

pub fn (mut s Server)accept() (&player.Player, bool) {
	h := s.server_accept(s.handle)
	if h == 0 {
		return unsafe {nil}, false
	}
	p := player.new(h) or {
		return unsafe {nil}, false
	}
	s.players[h] = p
	return p, true
}

pub fn (mut s Server)players() []&player.Player {
	mut players := []&player.Player{}
	pl := s.server_players(s.handle)
	arr := unsafe {arrays.carray_to_varray[voidptr](pl, int(sizeof(pl) / sizeof(voidptr)))}
	for p in arr {
		pi:= s.players[p] or {
			continue
		}
		players << pi
	}
	return players
}