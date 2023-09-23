package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"time"

	"github.com/faiface/pixel/pixelgl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/kbinani/screenshot"
	"github.com/otiai10/gosseract/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func main() {
	pixelgl.Run(run)
}

func run() {
	// Initialize GLFW
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Visible, glfw.False)
	window, err := glfw.CreateWindow(1, 1, "", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	// width, height := window.GetFramebufferSize()

	// Initialize Tesseract OCR client
	client := gosseract.NewClient()
	defer client.Close()

	time.Sleep(5 * time.Second)
	img, err := screenshot.CaptureRect(image.Rect(0, 0, 1920, 1080))
	if err != nil {
		fmt.Println("Error capturing screen:", err)
		return
	}
	err = saveImage(img, "captured.png")
	if err != nil {
		fmt.Println("Error saving captured image:", err)
		return
	}

	// image.Rect(0, 0, 1920, 1080)
	// take user input for the coordinates from the terminal
	fmt.Println("Enter the coordinates of the top left corner of the screen")
	var x1, y1 int
	fmt.Scanln(&x1, &y1)
	fmt.Println("Enter the coordinates of the bottom right corner of the screen")
	var x2, y2 int
	fmt.Scanln(&x2, &y2)

	for !window.ShouldClose() {

		// client.SetVariable("tessedit_char_whitelist", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

		// Capture the screen
		// width of my screen 0 to 1920
		// height of my screen 0 to 1080
		img, err := screenshot.CaptureRect(image.Rect(x1, y1, x2, y2))
		if err != nil {
			fmt.Println("Error capturing screen:", err)
			continue
		}

		fmt.Printf("Captured image resolution: %dx%d\n", img.Bounds().Dx(), img.Bounds().Dy())

		// err = saveResultImage(img, "captured.png")
		// if err != nil {
		// 	fmt.Println("Error saving captured image:", err)
		// 	return
		// }

		buf := new(bytes.Buffer)
		if err := png.Encode(buf, img); err != nil {
			fmt.Println("Error encoding captured image:", err)
			continue
		}

		// Use Tesseract to extract text
		client.SetLanguage("eng") // Set the language to English (adjust as needed)
		client.SetImageFromBytes(buf.Bytes())
		text, err := client.Text()
		if err != nil {
			fmt.Println("Error extracting text:", err)
			continue
		}

		// fmt.Println(text)
		// instead of text we can save it in a file
		f, err := os.Create("text.txt")
		if err != nil {
			fmt.Println(err)
			return
		}
		l, err := f.WriteString(text)
		if err != nil {
			fmt.Println(err)
			f.Close()
			return
		}
		fmt.Println(l, "bytes written successfully")
		err = f.Close()
		if err != nil {
			fmt.Println(err)
			return
		}

		// Delay between captures (adjust as needed)
		// glfw.PollEvents()
		time.Sleep(10 * time.Second)
	}
}

func saveResultImage(img *image.RGBA, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := png.Encode(file, img); err != nil {
		return err
	}
	return nil
}

func saveImage(img *image.RGBA, filename string) error {
	// addGrid(img)
	addGrid(img)

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := png.Encode(file, img); err != nil {
		return err
	}
	return nil
}

func addGrid(img *image.RGBA) {
	// Set the grid line color (e.g., red)
	gridColor := color.RGBA{255, 0, 0, 255}

	// Define the grid spacing (adjust as needed)
	gridSize := 100

	// Draw vertical grid lines
	for x := 0; x < img.Bounds().Dx(); x += gridSize {
		draw.Draw(img, image.Rect(x, 0, x+1, img.Bounds().Dy()), &image.Uniform{gridColor}, image.Point{}, draw.Over)
	}

	// Draw horizontal grid lines
	for y := 0; y < img.Bounds().Dy(); y += gridSize {
		draw.Draw(img, image.Rect(0, y, img.Bounds().Dx(), y+1), &image.Uniform{gridColor}, image.Point{}, draw.Over)
	}

	// Add labels to the grid (numbers)
	labelColor := color.RGBA{0, 0, 0, 255}
	labelFont := basicfont.Face7x13

	for x := 0; x < img.Bounds().Dx(); x += gridSize {
		for y := 0; y < img.Bounds().Dy(); y += gridSize {
			drawText(img, labelFont, x, y, fmt.Sprintf("(%d, %d)", x, y), labelColor)
		}
	}
}

func drawText(img *image.RGBA, face font.Face, x, y int, s string, col color.RGBA) {
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: face,
		Dot:  fixed.Point26_6{fixed.Int26_6(x << 6), fixed.Int26_6(y << 6)},
	}
	d.DrawString(s)
}
