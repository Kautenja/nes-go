
# Compile the C shared object library and Go wrapper in cascade
compile:
	cd nes && scons
	mv nes/lib_nes_env.so .
	mkdir -p build
	cd build && go build ../NESgo.go
