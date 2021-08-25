package main

// import (
// 	"github.com/gin-gonic/gin"
// 	"github.com/pion/rtcp"
// 	"github.com/pion/webrtc/v2"
// )

func main() {
    emulator := load("smb1.nes")
    reset(emulator)
    step(emulator)
    backup(emulator)
    restore(emulator)
    close(emulator)
}
