package main

import (
    "github.com/gorilla/websocket"
    "net/http"
    "log"
    "fmt"
    "encoding/base64"
    "time"
    "bytes"
    "image"
    "image/png"
);

func main() {
    http.HandleFunc("/screen/", ConnWs)
    err := http.ListenAndServe(":9090", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

func ConnWs(w http.ResponseWriter, r *http.Request) {
    ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
    _, ok := err.(websocket.HandshakeError)
    if ok {
        http.Error(w, "Not a websocket handshake", 400)
        return
    } else if err != nil {
        log.Println(err)
        return
    }

    emulator := load("smb1.nes")
    reset(emulator)

    res := map[string]interface{}{}
    	err = ws.ReadJSON(&res)
        if err != nil {
            if err.Error() == "EOF" {
                return
            }
            // ErrShortWrite means a write accepted fewer bytes than requested then failed to return an explicit error.
            if err.Error() == "unexpected EOF" {
                return
            }
            fmt.Println("Read : " + err.Error())
            return
        }

        res["a"] = "a"
        log.Println(res)

    // Create a placeholder image for passing pixels to the PNG encoder.
    img := image.NewRGBA(image.Rect(0, 0, screen_width(), screen_height()))
    // Loop infinitely to accept new connections.
    for {
        // Process a graphical frame on the emulator. This call blocks and is
        // relatively long running due to the number of CPU / PPU cycles per
        // frame.
	    step(emulator)
        // Update the pixels of the image that represents the screen.
		img.Pix = pixels(emulator)
        // Set the alpha channel to max value (they are 0 by default).
        // TODO: can this loop be replaced with an atomic switch that ignores
        // the alpha channel?
        // TODO: is there an atomic way to specify BGR instead of RGB?
        for i := 3; i < 240 * 256 * 4; i += 4 { img.Pix[i] = 255 }
        // Compress the screen into a PNG container.
		screenCompressed := new(bytes.Buffer)
		png.Encode(screenCompressed, img)
        // Convert the PNG image to a compressed base64 string to serve.
        str := base64.StdEncoding.EncodeToString(screenCompressed.Bytes())
        // Set the based64 image on the packet to send to the front-end
        res["img64"] = str
        err = ws.WriteJSON(&res)
        if err != nil {  // Handle the optional error
            fmt.Println("watch dir - Write : " + err.Error())
        }
        // Sleep to keep the server's tick-rate within NES specifications.
        time.Sleep(5 * time.Millisecond);
    }
}

