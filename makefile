
# Compile the C shared object library
compileC:
	cd nes && scons
	mv nes/lib_nes_env.so .

# Compile the GoLang application
compileGo:
	mkdir -p build
	cd build && go build ../*.go

# Compile the C shared object library and Go wrapper in cascade
compile: compileC compileGo

# Run the game engine server
run:
	./build/main

# Run an HTTP page server to serve pages
serve:
	python3 -m http.server 9000
