package main

import (
	"flag"
	"go-rpi-rgb-led-matrix/pkg/matrix"
	"go-rpi-rgb-led-matrix/tools"
	"image/color"
	"os"
	"os/signal"
	"syscall"
)

var (
	rows                   = flag.Int("led-rows", 32, "number of rows supported")
	cols                   = flag.Int("led-cols", 64, "number of columns supported")
	parallel               = flag.Int("led-parallel", 1, "number of daisy-chained panels")
	chain                  = flag.Int("led-chain", 2, "number of displays daisy-chained")
	brightness             = flag.Int("brightness", 100, "brightness (0-100)")
	hardwareMapping        = flag.String("led-gpio-mapping", "regular", "Name of GPIO mapping used.")
	showRefresh            = flag.Bool("led-show-refresh", false, "Show refresh rate.")
	inverseColors          = flag.Bool("led-inverse", false, "Switch if your matrix has inverse colors on.")
	disableHardwarePulsing = flag.Bool("led-no-hardware-pulse", false, "Don't use hardware pin-pulse generation.")
	pixelMapping           = flag.String("led-pixel-mapper", "U-mapper", "Pixel mapping from api")
)

func main() {
	config := &matrix.DefaultConfig
	config.Rows = *rows
	config.Cols = *cols
	config.Parallel = *parallel
	config.ChainLength = *chain
	config.Brightness = *brightness
	config.HardwareMapping = *hardwareMapping
	config.ShowRefreshRate = *showRefresh
	config.InverseColors = *inverseColors
	config.DisableHardwarePulsing = *disableHardwarePulsing
	config.PixelMapping = *pixelMapping

	m, err := matrix.NewRGBLedMatrix(config)
	fatal(err)

	tk := tools.NewToolKit(m)
	defer tk.Close()

	tk.DrawString("Cтр 1\nCтр 2\nCтр 3\nCтр 4\nCтр 5\n", color.RGBA{R: 255, G: 255, B: 255, A: 255}, nil)

	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)

	select {
	// Main routine sleeps until signal is received or global context is done.
	case <-osSignal:
		return
	}
}

func init() {
	flag.Parse()
}

func fatal(err error) {
	if err != nil {
		panic(err)
	}
}
