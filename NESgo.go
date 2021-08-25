package main

// To load the NES env package, point the linker to the compiled shared object
// library and provide a header file describing the C interface. The library
// name is "lib_nes_env", but the "lib" is implied by "-l" and must be omitted.
// "NESgo.h" provides a low-level C interface devoid of the notion of
// proprietary types, i.e., using `void*` for proprietary structures.

// #cgo LDFLAGS: -L${SRCDIR} -l_nes_env
// #include "NESgo.h"
import "C"
import "unsafe"

type Emulator struct {
	/// The path to the ROM file to load into the emulator
	path string
	/// A pointer to the emulator instance
	instance unsafe.Pointer
	/// The controller for player 1
	player1 *C.char
	/// The controller for player 2
	player2 *C.char
	/// The pointer to the underlying contiguous RAM buffer for querying pixels
	screen *C.char
	/// The screen buffer of 32-bit BGRx pixels
	pixels []byte
	/// the underlying RAM buffer
	ram *C.char
}

func load(path string) Emulator {
	// Initialize an emulator structure and assign the path.
	var emulator Emulator
	emulator.path = path
	// Create an NES emulator instance for the given ROM file and store the
	// pointer into the instance.
	emulator.instance = C.Initialize(C.CString(path))
	// Get pointers to the controller buffers.
	emulator.player1 = C.Controller(emulator.instance, 0)
	emulator.player2 = C.Controller(emulator.instance, 1)
	// Create a pointer to the underlying screen buffer.
	emulator.screen = C.Screen(emulator.instance)
	emulator.pixels = C.GoBytes(unsafe.Pointer(&emulator.screen), C.Width() * C.Height())
	// Get a reference to the RAM buffer.
	emulator.ram = C.Memory(emulator.instance)

	return emulator
}

func reset(emulator Emulator) {
	C.Reset(emulator.instance)
}

func step(emulator Emulator) {
	C.Step(emulator.instance)
}

func backup(emulator Emulator) {
	C.Backup(emulator.instance)
}

func restore(emulator Emulator) {
	C.Restore(emulator.instance)
}

func close(emulator Emulator) {
	C.Close(emulator.instance)
}

func main() {
	emulator := load("smb1.nes")
	reset(emulator)
	step(emulator)
	backup(emulator)
	restore(emulator)
	close(emulator)
}
