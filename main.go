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
        log.Fatal("Read : " + err.Error())
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
        for i := 0; i < 240 * 256 * 4; i += 4 {
            // Swap from BGR to RGB
            img.Pix[i + 0], img.Pix[i + 2] = img.Pix[i + 2], img.Pix[i + 0]
            // set alpha to max value
            img.Pix[i + 3] = 255
        }
        // Compress the screen into a PNG container.
		screenCompressed := new(bytes.Buffer)
		png.Encode(screenCompressed, img)
        // Convert the PNG image to a compressed base64 string to serve.
        str := base64.StdEncoding.EncodeToString(screenCompressed.Bytes())
        // Set the base64 image on the packet to send to the front-end
        response["img64"] = str
        err = ws.WriteJSON(&response)
        if err != nil {  // Handle the optional error
            log.Fatal("Screen Write : " + err.Error())
        }

        // Query the response from the server to update the controller.
        controllerResponse := map[string]interface{}{}
        err = ws.ReadJSON(&controllerResponse)
        if err != nil {
            if err.Error() == "EOF" {
                return
            }
            if err.Error() == "unexpected EOF" {
                return
            }
            log.Fatal("Read : " + err.Error())
            return
        }
        // Expect a float64 and convert to a byte to pass to the emulator.
        controller := byte(controllerResponse["controller"].(float64))
        player1(emulator, controller)

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
    if err != nil {  // Server failed to launch, report an error and terminate.
        log.Fatal("ListenAndServe: ", err)
    }
}
