package birch

import (
	"image"
	"image/color"
	"math"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/oakmound/oak/alg/floatgeom"
)

// A Triangle is a set of three points or vertices.
type Triangle struct {
	a, b, c mgl64.Vec3
}

// Unit returns a Triangle where each vertex is converted into a unit vector.
func (t Triangle) Unit() Triangle {
	return Triangle{Unit(t.a), Unit(t.b), Unit(t.c)}
}

// Mul returns a Triangle with each vertex magnified by f.
func (t Triangle) Mul(f float64) Triangle {
	return Triangle{t.a.Mul(f), t.b.Mul(f), t.c.Mul(f)}
}

// More documentation needed in the rest of this file

func (t Triangle) ViewTri(x, y, z, eye mgl64.Vec3) Triangle {
	return Triangle{
		mgl64.Vec3{t.a.Dot(x) - x.Dot(eye), t.a.Dot(y) - y.Dot(eye), t.a.Dot(z) - z.Dot(eye)},
		mgl64.Vec3{t.b.Dot(x) - x.Dot(eye), t.b.Dot(y) - y.Dot(eye), t.b.Dot(z) - z.Dot(eye)},
		mgl64.Vec3{t.c.Dot(x) - x.Dot(eye), t.c.Dot(y) - y.Dot(eye), t.c.Dot(z) - z.Dot(eye)},
	}
}

func (t Triangle) ViewNrm(x, y, z mgl64.Vec3) Triangle {
	return Triangle{
		mgl64.Vec3{t.a.Dot(x), t.a.Dot(y), t.a.Dot(z)},
		mgl64.Vec3{t.b.Dot(x), t.b.Dot(y), t.b.Dot(z)},
		mgl64.Vec3{t.c.Dot(x), t.c.Dot(y), t.c.Dot(z)},
	}.Unit()
}

func (t Triangle) Viewport(field floatgeom.Point2) Triangle {
	w := field.Y() / 1.5
	h := field.Y() / 1.5
	x := field.X() / 2.0
	y := field.Y() / 4.0
	return Triangle{
		mgl64.Vec3{w*t.a.X() + x, h*t.a.Y() + y, (t.a.Z() + 1.0) / 1.5},
		mgl64.Vec3{w*t.b.X() + x, h*t.b.Y() + y, (t.b.Z() + 1.0) / 1.5},
		mgl64.Vec3{w*t.c.X() + x, h*t.c.Y() + y, (t.c.Z() + 1.0) / 1.5},
	}
}

func (t Triangle) Perspective() Triangle {
	c := 3.0
	za := 1.0 - t.a.Z()/c
	zb := 1.0 - t.b.Z()/c
	zc := 1.0 - t.c.Z()/c
	return Triangle{
		mgl64.Vec3{t.a.X() / za, t.a.Y() / za, t.a.Z() / za},
		mgl64.Vec3{t.b.X() / zb, t.b.Y() / zb, t.b.Z() / zb},
		mgl64.Vec3{t.c.X() / zc, t.c.Y() / zc, t.c.Z() / zc},
	}
}

func (t Triangle) Translate(x, y float64) Triangle {
	t.a.Add(mgl64.Vec3{x, y})
	t.b.Add(mgl64.Vec3{x, y})
	t.c.Add(mgl64.Vec3{x, y})

	return Triangle{
		a: t.a,
		b: t.b,
		c: t.c,
	}
}

func (t Triangle) BaryCenter(x, y int) mgl64.Vec3 {
	p := mgl64.Vec3{float64(x), float64(y), 0.0}
	v0 := t.b.Sub(t.a)
	v1 := t.c.Sub(t.a)
	v2 := p.Sub(t.a)
	d00 := v0.Dot(v0)
	d01 := v0.Dot(v1)
	d11 := v1.Dot(v1)
	d20 := v2.Dot(v0)
	d21 := v2.Dot(v1)
	v := (d11*d20 - d01*d21) / (d00*d11 - d01*d01)
	w := (d00*d21 - d01*d20) / (d00*d11 - d01*d01)
	u := 1.0 - v - w
	return mgl64.Vec3{v, w, u}
}

func TDraw(buff *image.RGBA, zbuff [][]float64, vew, nrm, tex Triangle, textureData *image.RGBA) {
	x0 := int(math.Min(vew.a.X(), math.Min(vew.b.X(), vew.c.X())))
	y0 := int(math.Min(vew.a.Y(), math.Min(vew.b.Y(), vew.c.Y())))
	x1 := int(math.Max(vew.a.X(), math.Max(vew.b.X(), vew.c.X())))
	y1 := int(math.Max(vew.a.Y(), math.Max(vew.b.Y(), vew.c.Y())))
	dims := textureData.Bounds()
	buffH := buff.Bounds().Max.Y
	for x := x0; x <= x1; x++ {
		for y := y0; y <= y1; y++ {
			bc := vew.BaryCenter(x, y)
			if bc.X() >= 0.0 && bc.Y() >= 0.0 && bc.Z() >= 0.0 {
				// Multiply everything by Z to create perspective.
				z := bc.X()*vew.b.Z() + bc.Y()*vew.c.Z() + bc.Z()*vew.a.Z()

				if z > zbuff[x][y] {
					light := mgl64.Vec3{0.0, 0.0, 1.0}
					varying := mgl64.Vec3{light.Dot(nrm.b), light.Dot(nrm.c), light.Dot(nrm.a)}

					xx := (float64(dims.Max.X) - 1) * (0.0 + (bc.X()*tex.b.X() + bc.Y()*tex.c.X() + bc.Z()*tex.a.X()))
					yy := (float64(dims.Max.Y) - 1) * (1.0 - (bc.X()*tex.b.Y() + bc.Y()*tex.c.Y() + bc.Z()*tex.a.Y()))
					intensity := bc.Dot(varying)
					var shading uint32
					if intensity > 0.0 {
						shading = uint32(intensity * 0xFF)
					}
					zbuff[x][y] = z
					// Change from the original gel: we subtract y from buffH because,
					// somewhere, I (200sc) messed up the translation and we accidentally
					// are rendering everything upsidedown. This is the easiest fix!
					buff.Set(x, buffH-y, PShade(textureData.At(int(xx), int(yy)), shading))
				}
			}
		}
	}
}

// PShade takes a color and applies shading to it by magnifying it's rgb values
func PShade(pixel color.Color, shading uint32) color.RGBA {
	r, g, b, a := pixel.RGBA()
	// r,g, and b are divided by 257 because the .RGBA() function returns
	// values on a 16 bit scale instead of an 8 bit scale.
	// They are then magnified by shading, not overflowing because they are
	// int32 values, and shifted to the right to represent division after
	// shading.
	// Todo: shading should probably by a float64 instead of multiplying and
	// then shifting
	r = ((r / 257) * shading) >> 0x08
	g = ((g / 257) * shading) >> 0x08
	b = ((b / 257) * shading) >> 0x08
	// r,g, and b need to be cast back to uint8s to create a color.RGBA value
	// from them.
	return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}
