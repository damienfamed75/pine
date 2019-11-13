// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This build tag means that "go install golang.org/x/exp/shiny/..." doesn't
// install this example program. Use "go run main.go" to run it or "go install
// -tags=example" to install it.

// Basic is a basic example of a graphical application.
package main

import (
	"image"
	"image/color"
	"image/draw"
	"log"
	"math"

	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/shiny/driver"
	"github.com/oakmound/shiny/screen"
	"github.com/oakmound/shiny/widget/glwidget"
	"github.com/oakmound/shiny/widget/node"
	"github.com/oakmound/shiny/widget/theme"
)

var (
	blue0    = color.RGBA{0x00, 0x00, 0x1f, 0xff}
	blue1    = color.RGBA{0x00, 0x00, 0x3f, 0xff}
	darkGray = color.RGBA{0x3f, 0x3f, 0x3f, 0xff}
	green    = color.RGBA{0x00, 0x7f, 0x00, 0x7f}
	red      = color.RGBA{0x7f, 0x00, 0x00, 0x7f}
	yellow   = color.RGBA{0x3f, 0x3f, 0x00, 0x3f}

	cos30 = math.Cos(math.Pi / 6)
	sin30 = math.Sin(math.Pi / 6)

	scrWidth  = 800
	scrHeight = 600

	screenControl screen.Screen
	windowControl screen.Window
	winBuffer     screen.Image

	windowRect image.Rectangle
)

func main() {
	driver.Main(func(s screen.Screen) {
		screenControl = s

		var err error
		winBuffer, err = screenControl.NewImage(image.Point{scrWidth, scrHeight})
		if err != nil {
			log.Fatalf("error creating winbuffer: %s", err)
		}
		glwidget.NewGL(drawer)

		changeWindow(0, 0, scrWidth, scrHeight)

		for {

			tx, err := screenControl.NewTexture(winBuffer.Bounds().Max)
			if err != nil {
				log.Fatalf("failed to create win texture: %s", err)
			}

			tx.Upload(image.Point{0, 0}, winBuffer, winBuffer.Bounds())
			windowControl.Scale(windowRect, tx, tx.Bounds(), draw.Src)
			windowControl.Publish()
			// winBuffer.RGBA().SetRGBA(10, 10, blue0)
			// winBuffer.RGBA().SetRGBA(10, 11, blue0)
			// winBuffer.RGBA().SetRGBA(10, 12, blue0)
			// winBuffer.RGBA().SetRGBA(10, 13, blue0)

		}

	})

	// scr := screen.S(screen.Title("Hello,World"))

	// win := screen.NewWindowGenerator(screen.Dimensions(800,600))

	// winBuff := scr.
}

func drawer(glctx *glwidget.GL) {
	ctx := node.PaintBaseContext{
		Theme: theme.Default,
		Dst:   winBuffer.RGBA(),
	}
	glctx.PaintBase(&ctx, image.Point{0, 0})
}

func windowController(s screen.Screen, x, y int32, width, height int) (screen.Window, error) {
	return s.NewWindow(screen.NewWindowGenerator(
		screen.Dimensions(width, height),
		screen.Title("Hello, World!"),
		screen.Position(x, y),
		screen.Fullscreen(false),
		screen.Borderless(false),
		screen.TopMost(false),
	))
}

func changeWindow(x, y int32, width, height int) {
	// The window controller handles incoming hardware or platform events and
	// publishes image data to the screen.
	wC, err := windowController(screenControl, x, y, width, height)
	if err != nil {
		dlog.Error(err)
		panic(err)
	}
	windowControl = wC
	ChangeWindow(width, height)
}

// ChangeWindow sets the width and height of the game window. Although exported,
// calling it without a size event will probably not act as expected.
func ChangeWindow(width, height int) {
	// Draw a black frame to cover up smears
	// Todo: could restrict the black to -just- the area not covered by the
	// scaled screen buffer
	buff, err := screenControl.NewImage(image.Point{width, height})
	if err == nil {
		draw.Draw(buff.RGBA(), buff.Bounds(), image.Black, image.Point{0, 0}, draw.Src)
		windowControl.Upload(image.Point{0, 0}, buff, buff.Bounds())
	} else {
		dlog.Error(err)
	}
	var x, y int
	// if UseAspectRatio {
	// 	inRatio := float64(width) / float64(height)
	// 	if aspectRatio > inRatio {
	// 		newHeight := alg.RoundF64(float64(height) * (inRatio / aspectRatio))
	// 		y = (newHeight - height) / 2
	// 		height = newHeight - y
	// 	} else {
	// 		newWidth := alg.RoundF64(float64(width) * (aspectRatio / inRatio))
	// 		x = (newWidth - width) / 2
	// 		width = newWidth - x
	// 	}
	// }
	windowRect = image.Rect(-x, -y, width, height)
}

// //+build darwin,metal

// package main

// import (
// 	"encoding/binary"
// 	"fmt"
// 	"image"
// 	"image/color"
// 	"log"
// 	"math"

// 	// "golang.org/x/exp/shiny/driver/gldriver"
// 	// "golang.org/x/exp/shiny/screen"
// 	// "golang.org/x/exp/shiny/unit"
// 	// "golang.org/x/exp/shiny/widget"
// 	// "golang.org/x/exp/shiny/widget/flex"
// 	// "golang.org/x/exp/shiny/widget/glwidget"
// 	// "golang.org/x/exp/shiny/widget/theme"
// 	"github.com/oakmound/shiny/driver/gldriver"
// 	"github.com/oakmound/shiny/screen"
// 	"github.com/oakmound/shiny/unit"
// 	"github.com/oakmound/shiny/widget"
// 	"github.com/oakmound/shiny/widget/flex"
// 	"github.com/oakmound/shiny/widget/glwidget"
// 	"github.com/oakmound/shiny/widget/theme"
// 	"golang.org/x/image/colornames"
// 	"golang.org/x/mobile/gl"
// )

// const (
// 	screenWidth  = 800
// 	screenHeight = 600
// )

// func colorPatch(c color.Color, w, h unit.Value) *widget.Sizer {
// 	return widget.NewSizer(w, h, widget.NewUniform(theme.StaticColor(c), nil))
// }

// func main() {
// 	gldriver.Main(func(s screen.Screen) {
// 		t1, t2 := newTriangleGL(), newTriangleGL()
// 		defer t1.cleanup() // free the memory
// 		defer t2.cleanup() // free the memory

// 		body := widget.NewSheet(flex.NewFlex(
// 			colorPatch(colornames.Green, unit.Pixels(50), unit.Pixels(50)),
// 			widget.WithLayoutData(t1.w, flex.LayoutData{Grow: 1, Align: flex.AlignItemStretch}),
// 			colorPatch(colornames.Blue, unit.Pixels(50), unit.Pixels(50)),
// 			widget.WithLayoutData(t2.w, flex.LayoutData{MinSize: image.Point{80, 80}}),
// 			colorPatch(colornames.Green, unit.Pixels(50), unit.Pixels(50)),
// 		))

// 		if err := widget.RunWindow(s, body, &widget.RunWindowOptions{
// 			// NewWindowOptions: screen.WindowGenerator{
// 			// 	Title: "Hello, World!",
// 			// },
// 			// NewWindowOptions: screen.NewWindowOptions{
// 			// 	Title: "Hello, World!",
// 			// },
// 		}); err != nil {
// 			log.Fatalf("failed to run window: %s", err)
// 		}
// 	})
// }

// func newTriangleGL() *triangleGL {
// 	t := &triangleGL{}
// 	t.w = glwidget.NewGL(t.draw)
// 	t.init()
// 	return t
// }

// type triangleGL struct {
// 	w *glwidget.GL

// 	program  gl.Program
// 	position gl.Attrib
// 	offset   gl.Uniform
// 	color    gl.Uniform
// 	buf      gl.Buffer

// 	green float32
// }

// func (t *triangleGL) init() {
// 	glctx := t.w.Ctx
// 	var err error
// 	t.program, err = createProgram(glctx, vertexShader, fragmentShader)
// 	if err != nil {
// 		log.Fatalf("error creating GL program: %v", err)
// 	}

// 	t.buf = glctx.CreateBuffer()
// 	glctx.BindBuffer(gl.ARRAY_BUFFER, t.buf)
// 	glctx.BufferData(gl.ARRAY_BUFFER, triangleData, gl.STATIC_DRAW)

// 	t.position = glctx.GetAttribLocation(t.program, "position")
// 	t.color = glctx.GetUniformLocation(t.program, "color")
// 	t.offset = glctx.GetUniformLocation(t.program, "offset")

// 	glctx.UseProgram(t.program)
// 	glctx.ClearColor(1, 0, 0, 1)
// }

// func (t *triangleGL) cleanup() {
// 	glctx := t.w.Ctx
// 	glctx.DeleteProgram(t.program)
// 	glctx.DeleteBuffer(t.buf)
// }

// func (t *triangleGL) draw(w *glwidget.GL) {
// 	glctx := t.w.Ctx

// 	glctx.Viewport(0, 0, w.Rect.Dx(), w.Rect.Dy())
// 	glctx.Clear(gl.COLOR_BUFFER_BIT)

// 	t.green += 0.01
// 	if t.green > 1 {
// 		t.green = 0
// 	}
// 	glctx.Uniform4f(t.color, 0, t.green, 0, 1)
// 	glctx.Uniform2f(t.offset, 0.2, 0.9)

// 	glctx.BindBuffer(gl.ARRAY_BUFFER, t.buf)
// 	glctx.EnableVertexAttribArray(t.position)
// 	glctx.VertexAttribPointer(t.position, coordsPerVertex, gl.FLOAT, false, 0, 0)
// 	glctx.DrawArrays(gl.TRIANGLES, 0, vertexCount)
// 	glctx.DisableVertexAttribArray(t.position)
// 	w.Publish()
// }

// // asBytes returns the byte representation of float32 values in the given byte
// // order. byteOrder must be either binary.BigEndian or binary.LittleEndian.
// func asBytes(byteOrder binary.ByteOrder, values ...float32) []byte {
// 	le := false
// 	switch byteOrder {
// 	case binary.BigEndian:
// 	case binary.LittleEndian:
// 		le = true
// 	default:
// 		panic(fmt.Sprintf("invalid byte order %v", byteOrder))
// 	}

// 	b := make([]byte, 4*len(values))
// 	for i, v := range values {
// 		u := math.Float32bits(v)
// 		if le {
// 			b[4*i+0] = byte(u >> 0)
// 			b[4*i+1] = byte(u >> 8)
// 			b[4*i+2] = byte(u >> 16)
// 			b[4*i+3] = byte(u >> 24)
// 		} else {
// 			b[4*i+0] = byte(u >> 24)
// 			b[4*i+1] = byte(u >> 16)
// 			b[4*i+2] = byte(u >> 8)
// 			b[4*i+3] = byte(u >> 0)
// 		}
// 	}
// 	return b
// }

// // createProgram creates, compiles, and links a gl.Program.
// func createProgram(glctx gl.Context, vertexSrc, fragmentSrc string) (gl.Program, error) {
// 	program := glctx.CreateProgram()
// 	if program.Value == 0 {
// 		return gl.Program{}, fmt.Errorf("basicgl: no programs available")
// 	}

// 	vertexShader, err := loadShader(glctx, gl.VERTEX_SHADER, vertexSrc)
// 	if err != nil {
// 		return gl.Program{}, err
// 	}
// 	fragmentShader, err := loadShader(glctx, gl.FRAGMENT_SHADER, fragmentSrc)
// 	if err != nil {
// 		glctx.DeleteShader(vertexShader)
// 		return gl.Program{}, err
// 	}

// 	glctx.AttachShader(program, vertexShader)
// 	glctx.AttachShader(program, fragmentShader)
// 	glctx.LinkProgram(program)

// 	// Flag shaders for deletion when program is unlinked.
// 	glctx.DeleteShader(vertexShader)
// 	glctx.DeleteShader(fragmentShader)

// 	if glctx.GetProgrami(program, gl.LINK_STATUS) == 0 {
// 		defer glctx.DeleteProgram(program)
// 		return gl.Program{}, fmt.Errorf("basicgl: %s", glctx.GetProgramInfoLog(program))
// 	}
// 	return program, nil
// }

// func loadShader(glctx gl.Context, shaderType gl.Enum, src string) (gl.Shader, error) {
// 	shader := glctx.CreateShader(shaderType)
// 	if shader.Value == 0 {
// 		return gl.Shader{}, fmt.Errorf("basicgl: could not create shader (type %v)", shaderType)
// 	}
// 	glctx.ShaderSource(shader, src)
// 	glctx.CompileShader(shader)
// 	if glctx.GetShaderi(shader, gl.COMPILE_STATUS) == 0 {
// 		defer glctx.DeleteShader(shader)
// 		return gl.Shader{}, fmt.Errorf("basicgl: shader compile: %s", glctx.GetShaderInfoLog(shader))
// 	}
// 	return shader, nil
// }

// var triangleData = asBytes(binary.LittleEndian,
// 	0.0, 0.4, 0.0, // top left
// 	0.0, 0.0, 0.0, // bottom left
// 	0.4, 0.0, 0.0, // bottom right
// )

// const (
// 	coordsPerVertex = 3
// 	vertexCount     = 3
// )

// const vertexShader = `#version 100
// uniform vec2 offset;
// attribute vec4 position;
// void main() {
// 	// offset comes in with x/y values between 0 and 1.
// 	// position bounds are -1 to 1.
// 	vec4 offset4 = vec4(2.0*offset.x-1.0, 1.0-2.0*offset.y, 0, 0);
// 	gl_Position = position + offset4;
// }`

// const fragmentShader = `#version 100
// precision mediump float;
// uniform vec4 color;
// void main() {
// 	gl_FragColor = color;
// }`
