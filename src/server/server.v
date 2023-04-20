module server

import dl.loader
import player

type T_server_start = fn() voidptr
type T_server_accept = fn(voidptr) voidptr

pub struct Server{
	handle voidptr
	server_accept T_server_accept

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
	return p, true
}