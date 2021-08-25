
# Compile the C shared object library and Go wrapper in cascade
compile:
	cd nes && scons
	mv nes/lib_nes_env.so .
	go build NESgo.go
