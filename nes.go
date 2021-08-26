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

// @brief An NES emulator instance.
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
    /// the underlying RAM buffer
    ram *C.char
}

// @brief Load a ROM at the given path.
// @param path the path to the ROM file to load
// @returns a new emulator instance with the loaded ROM
//
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
    // emulator.pixels = C.GoBytes(unsafe.Pointer(&emulator.screen), 4 * C.Width() * C.Height())
    // Get a reference to the RAM buffer.
    emulator.ram = C.Memory(emulator.instance)
    // Reset the emulator to start the emulation.
    C.Reset(emulator.instance)

    return emulator
}

// @brief Return the number of vertical pixels on the screen.
func screen_height() int { return int(C.Height()) }

// @brief Return the number of horizontal pixels on the screen.
func screen_width() int { return int(C.Width()) }

// @brief Return the pixels of the screen.
// @param emulator the emulator to get the pixels
//
func pixels(emulator Emulator) []byte {
   return C.GoBytes(unsafe.Pointer(emulator.screen), 4 * C.Width() * C.Height())
}

// @brief Update the controller for player 1.
// @param emulator the emulator to update the controller of
// @param controller the bitmap representation of the buttons that are pressed
//
func player1(emulator Emulator, controller byte) {
    *emulator.player1 = C.char(controller)
}

// @brief Update the controller for player 2.
// @param emulator the emulator to update the controller of
// @param controller the bitmap representation of the buttons that are pressed
//
func player2(emulator Emulator, controller byte) {
    *emulator.player2 = C.char(controller)
}

// @brief Reset the emulator, i.e., like hitting the reset button on the NES.
// @param emulator the emulator to reset
//
func reset(emulator Emulator) { C.Reset(emulator.instance) }

// @brief Step the emulator forward a single video frame.
// @param emulator the emulator to step
//
func step(emulator Emulator) { C.Step(emulator.instance) }

// @brief Backup the state of the emulator.
// @param emulator the emulator to backup
//
func backup(emulator Emulator) { C.Backup(emulator.instance) }

// @brief Restore the state of the emulator from a backup.
// @param emulator the emulator to restore
//
func restore(emulator Emulator) { C.Restore(emulator.instance) }

// @brief Close the emulator. The struct is deferred past this point.
// @param emulator the emulator to close
//
func close(emulator Emulator) { C.Close(emulator.instance) }
