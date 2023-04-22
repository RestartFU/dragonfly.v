run:
	go build -o dragonfly.dll -buildmode=c-shared ./lib
	del dragonfly.h
	v -cc gcc run .