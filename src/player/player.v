module player

import dl.loader

type T_player_message = fn(voidptr, &u8)
type T_player_name = fn(voidptr) &u8
type T_player_handle_quit = fn(voidptr, fn())

pub struct Player {
	handle voidptr

	player_message T_player_message
	player_name T_player_name
	player_handle_quit T_player_handle_quit
}

pub fn new(h voidptr) !&Player {
	mut libs := []string{}
	libs << "dragonfly.dll"

	mut lib := loader.get_or_create_dynamic_lib_loader(loader.DynamicLibLoaderConfig{paths: libs}) or {
		panic(err)
	}

	lib.open()!

	p := &Player{
		handle: h
		player_message: T_player_message(lib.get_sym("player_message")!)
		player_name: T_player_name(lib.get_sym("player_name")!)
		player_handle_quit: T_player_handle_quit(lib.get_sym("player_handle_quit")!) // TODO: fix: unexpected fault address 0xffffffffffffffff
	}
	//p.player_handle_quit(p.handle, fn () {
	//	println("test") //
	//})
	return p
}

pub fn (mut p Player)message(msg string) {
	p.player_message(p.handle, msg.str)
}

pub fn (mut p Player)name() string{
	name := p.player_name(p.handle)
	return unsafe {cstring_to_vstring(name)}
}