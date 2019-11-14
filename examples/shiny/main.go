package main

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/exp/app/debug"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/gl"
)

var (
	touchX       float32
	touchY       float32
	green        float32 // Color of the triangle.
	triangleData = f32.Bytes(binary.LittleEndian,
		0.0, 0.4, 0.0, // top left
		0.0, 0.0, 0.0, // bottom left
		0.4, 0.0, 0.0, // bottom right
	)
)

const (
	// Number of vertices in the triangle
	vertexCount     = 3
	coordsPerVertex = 3

	scrWidth  = 800
	scrHeight = 600
)

type Context struct {
	window *screen.Window
	// Graphics Library Context.
	glctx gl.Context
	// Size after resizing window.
	sz size.Event
	// Shared info among the image objects.
	images *glutil.Images
	// Shader programs.
	programs map[string]ShaderProgram
	buffers  []gl.Buffer
	fps      *debug.FPS
}

type ShaderProgram struct {
	// The compiled shader program.
	executable gl.Program
	buffer     gl.Buffer
	attribs    map[string]gl.Attrib
	uniforms   map[string]gl.Uniform
}

func main() {
	driver.Main(func(s screen.Screen) {
		window, err := s.NewWindow(&screen.NewWindowOptions{
			Width:  800,
			Height: 600,
			Title:  "Hello, World!",
		})
		if err != nil {
			log.Fatalf("creating window: %s", err)
		}
		quit := make(chan struct{})

		ctx := &Context{
			window:   &window,
			programs: make(map[string]ShaderProgram),
		}

		for {
			// For every incoming event.
			switch e := window.NextEvent().(type) {
			case lifecycle.Event: // General life of the window.
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn: // Begin window
					log.Printf("Intializing Window")
					ctx.glctx, _ = e.DrawContext.(gl.Context)

					onStart(ctx)
					// Begin drawing loop.
					// window.Send(paint.Event{})
					go func() {
						for {
							select {
							case <-quit:
								return
							default:
							}
							if ctx.glctx == nil {
								continue
							}

							onPaint(ctx)
							window.Publish()
						}
					}()
				case lifecycle.CrossOff: // Stop window
					log.Printf("Killing Window")
					onStop(ctx)
					close(quit)
					ctx.glctx = nil
					return
				}
			case size.Event: // Resizing window.
				ctx.sz = e
				touchX = float32(ctx.sz.WidthPx / 2)
				touchY = float32(ctx.sz.HeightPx / 2)
			case paint.Event: // Painting to the window.
				// paint...
				// if ctx.glctx == nil || e.External {
				// 	continue
				// }
				// onPaint(ctx)
				// window.Publish()
				// window.Send(paint.Event{})
			case key.Event:
				if e.Code == key.CodeEscape {
					window.Send(lifecycle.Event{
						From: lifecycle.StageVisible,
						To:   lifecycle.StageDead,
					})
				}
			}
		}
	})
}

func onStart(ctx *Context) error {
	vertShader, err := loadShader("basic.vert")
	if err != nil {
		return fmt.Errorf("loading vertShader: %w", err)
	}

	fragShader, err := loadShader("basic.frag")
	if err != nil {
		return fmt.Errorf("loading fragShader: %w", err)
	}

	program, err := glutil.CreateProgram(ctx.glctx, vertShader, fragShader)
	if err != nil {
		return fmt.Errorf("creating shader program: %w", err)
	}

	buffer := ctx.glctx.CreateBuffer()
	ctx.glctx.BindBuffer(gl.ARRAY_BUFFER, buffer)
	ctx.glctx.BufferData(gl.ARRAY_BUFFER, triangleData, gl.STATIC_DRAW)
	// Append the newly created buffer object into the context's buffer list.
	ctx.buffers = append(ctx.buffers, buffer)

	position := ctx.glctx.GetAttribLocation(program, "position")
	color := ctx.glctx.GetUniformLocation(program, "color")
	offset := ctx.glctx.GetUniformLocation(program, "offset")

	// Append the final program to the list of shader programs.
	ctx.programs["basic"] = ShaderProgram{
		executable: program,
		buffer:     buffer,
		attribs: map[string]gl.Attrib{
			"position": position,
		},
		uniforms: map[string]gl.Uniform{
			"color":  color,
			"offset": offset,
		},
	}

	ctx.images = glutil.NewImages(ctx.glctx)
	ctx.fps = debug.NewFPS(ctx.images)

	return nil
}

// onStop releases all the memory allocated for the application.
func onStop(ctx *Context) {
	for _, p := range ctx.programs {
		ctx.glctx.DeleteProgram(p.executable)
	}

	for _, b := range ctx.buffers {
		ctx.glctx.DeleteBuffer(b)
	}

	ctx.fps.Release()
	ctx.images.Release()
}

func onPaint(ctx *Context) {
	ctx.glctx.ClearColor(1, 0, 0, 1)
	ctx.glctx.Clear(gl.COLOR_BUFFER_BIT)

	green += 0.01
	if green > 1 {
		green = 0
	}

	// Load the basic shader.
	shader := ctx.programs["basic"]
	ctx.glctx.UseProgram(shader.executable)
	// Set the color value.
	// (R, G, B, A)
	ctx.glctx.Uniform4f(
		shader.uniforms["color"],
		0, green, 0, 1,
	)
	// Set the offset.
	// (X, Y)
	ctx.glctx.Uniform2f(
		shader.uniforms["offset"],
		touchX/float32(ctx.sz.WidthPx),
		touchY/float32(ctx.sz.HeightPx),
	)

	// Bind the triangle data.
	// This includes the vertices of the triangle.
	ctx.glctx.BindBuffer(gl.ARRAY_BUFFER, shader.buffer)

	// Get the uid of the position variable.
	position := shader.attribs["position"]

	// Render at position.
	ctx.glctx.EnableVertexAttribArray(position)
	{
		// With the given bound buffer to ARRAY_BUFFER
		// VertexAttribPointer binds to a generic vertex attribute of the
		// current vertex buffer object and specifies its layout.
		ctx.glctx.VertexAttribPointer(
			position,
			// Number of elements per vertex attribute.
			// For instance a triangle would be 3 for three vertices.
			coordsPerVertex,
			// Datatype
			gl.FLOAT,
			// Normalized
			false,
			// Stride is the byte offset between the vertex components.
			0,
			// Offset in bytes of the first component in the vertex attribute
			// array. Must be a multiple of the byte length of the datatype.
			0,
		)
		// Draw the buffer.
		// (mode, first, count)
		ctx.glctx.DrawArrays(gl.TRIANGLES, 0, vertexCount)
	}
	// Stops rendering at the given position.
	// (attrib uid)
	ctx.glctx.DisableVertexAttribArray(position)

	// Draw the FPS counter.
	ctx.fps.Draw(ctx.sz)
}

func loadShader(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	bVal, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(bVal), nil
}
