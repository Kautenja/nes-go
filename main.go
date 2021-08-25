package main

import (
    "github.com/gorilla/websocket"
    "net/http"
    "log"
    "encoding/base64"
    "time"
    "bytes"
    "image"
    "image/png"
);

// @brief Handle a request for the screen endpoint
// @param writer the HTTP response writer for sending data
// @param request the HTTP request that was received at the endpoint
//
func screen(writer http.ResponseWriter, request *http.Request) {
    ws, err := websocket.Upgrade(writer, request, nil, 1024, 1024)
    _, ok := err.(websocket.HandshakeError)
    if ok {
        http.Error(writer, "This end-point expects a WebSocket handshake", 400)
        return
    } else if err != nil {
        log.Println(err)
        return
    }

    // Create an emulator instance for this connection.
    emulator := load("smb1.nes")

    response := map[string]interface{}{}
    err = ws.ReadJSON(&response)
    if err != nil {
        if err.Error() == "EOF" {
            return
        }
        // ErrShortWrite means a write accepted fewer bytes than requested
        // then failed to return an explicit error.
        if err.Error() == "unexpected EOF" {
            return
        }
        print("Read : " + err.Error())
        return
    }

    response["a"] = "a"
    log.Println(response)

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
        response["img64"] = str
        err = ws.WriteJSON(&response)
        if err != nil {  // Handle the optional error
            print("watch dir - Write : " + err.Error())
        }
        // Sleep to keep the server's tick-rate within NES specifications. The
        // NES ran at 60Hz = 16.7ms, but there is some overhead associated with
        // the network stack. 5ms works well in practice, but this will need to
        // be tuned / refactored to lock the tick rate to 60Hz properly.
        time.Sleep(5 * time.Millisecond);
    }
}

// @brief The main entry point.
func main() {
    // Setup the functional callbacks for the endpoints
    http.HandleFunc("/screen/", screen)
    // Start the server and handle any error that occurs
    err := http.ListenAndServe(":9090", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
