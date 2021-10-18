package main

import (
	"flag"
	"fmt"
	"github.com/BoskyWSMFN/go-rpi-rgb-led-matrix/tools"
	"os"
	"time"

	"github.com/BoskyWSMFN/go-rpi-rgb-led-matrix/pkg/rpc"
)

var (
	img = flag.String("image", "", "image path")
)

func main() {
	f, err := os.Open(*img)
	fatal(err)

	m, err := rpc.NewClient("tcp", "10.42.0.161:1234", 100*100)
	fatal(err)

	tk := tools.NewToolKit(m)
	close, err := tk.PlayGIF(f)
	fatal(err)

	time.Sleep(time.Second * 3)
	close <- true
}

func init() {
	flag.Parse()
}

func fatal(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
