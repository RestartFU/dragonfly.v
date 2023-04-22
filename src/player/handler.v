module player

interface Handler {
mut:
	handle_quit()
}

struct NopHandler {}

fn (mut _ NopHandler) handle_quit() {}
