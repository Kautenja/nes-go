package main

func main() {
    emulator := load("smb1.nes")
    reset(emulator)
    step(emulator)
    backup(emulator)
    restore(emulator)
    close(emulator)
}
