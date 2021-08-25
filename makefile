
# Compile the C shared object library and Go wrapper in cascade
compileC:
	cd nes && scons
	mv nes/lib_nes_env.so .

compileGo:
	mkdir -p build
	cd build && go build ../*.go

compile: compileC compileGo
