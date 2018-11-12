package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	astilectron "github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bundler"
)

func main() {
	var a, _ = astilectron.New(astilectron.Options{
		AppName:            "hashira",
		AppIconDarwinPath:  "resources/icon.icns",
		AppIconDefaultPath: "resources/icon.png",
	})
	defer a.Close()

	// Set provisioner
	a.SetProvisioner(astibundler.NewProvisioner(Asset))

	rp := filepath.Join(a.Paths().BaseDirectory(), "resources")
	_, err := os.Stat(rp)
	if os.IsNotExist(err) {
		RestoreAssets(a.Paths().BaseDirectory(), "resources")
	}

	// Start astilectron
	a.Start()

	// Create a new window
	var w, _ = a.NewWindow(filepath.Join(a.Paths().BaseDirectory(), "resources", "app", "index.html"),
		&astilectron.WindowOptions{
			Center: astilectron.PtrBool(true),
			Height: astilectron.PtrInt(600),
			Width:  astilectron.PtrInt(600),
		})
	w.Create()
	w.OnMessage(handleMessage)

	// TODO: enable only for debug
	w.OpenDevTools()
	a.Wait()
}

type message struct {
	Name    string          `json:"name"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

func handleMessage(m *astilectron.EventMessage) interface{} {
	bytes, err := m.MarshalJSON()
	if err != nil {
		fmt.Println("MarshalJSON() error. err =", err.Error())
		return nil
	}
	fmt.Println(string(bytes))

	var in message
	err = m.Unmarshal(&in)
	if err != nil {
		fmt.Println("UnmarshalJSON() error. err =", err.Error())
		return nil
	}
	fmt.Println(in.Name)
	fmt.Println(string(in.Payload))

	return &message{Name: in.Name + ".callback", Payload: []byte("hoge")}
}
