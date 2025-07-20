package primitive

import (
	"fmt"
	"math"

	"github.com/fogleman/gg"
)

type Triangle struct {
	Worker *Worker
	X1, Y1 int
	X2, Y2 int
	X3, Y3 int
}

func NewRandomTriangle(worker *Worker) *Triangle {
	rnd := worker.Rnd
	x1 := rnd.Intn(worker.W)
	y1 := rnd.Intn(worker.H)
	x2 := x1 + rnd.Intn(31) - 15
	y2 := y1 + rnd.Intn(31) - 15
	x3 := x1 + rnd.Intn(31) - 15
	y3 := y1 + rnd.Intn(31) - 15
	t := &Triangle{worker, x1, y1, x2, y2, x3, y3}
	t.Mutate()
	return t
}

func (t *Triangle) Draw(dc *gg.Context, scale float64) {
	dc.LineTo(float64(t.X1), float64(t.Y1))
	dc.LineTo(float64(t.X2), float64(t.Y2))
	dc.LineTo(float64(t.X3), float64(t.Y3))
	dc.ClosePath()
	dc.Fill()
}

func (t *Triangle) SVG(attrs string) string {
	return fmt.Sprintf(
		"<polygon %s points=\"%d,%d %d,%d %d,%d\" />",
		attrs, t.X1, t.Y1, t.X2, t.Y2, t.X3, t.Y3)
}

func (t *Triangle) Copy() Shape {
	a := *t
	return &a
}

func (t *Triangle) Mutate() {
	w := t.Worker.W
	h := t.Worker.H
	rnd := t.Worker.Rnd
	const m = 16
	for {
		switch rnd.Intn(3) {
		case 0:
			t.X1 = clampInt(t.X1+int(rnd.NormFloat64()*16), -m, w-1+m)
			t.Y1 = clampInt(t.Y1+int(rnd.NormFloat64()*16), -m, h-1+m)
		case 1:
			t.X2 = clampInt(t.X2+int(rnd.NormFloat64()*16), -m, w-1+m)
			t.Y2 = clampInt(t.Y2+int(rnd.NormFloat64()*16), -m, h-1+m)
		case 2:
			t.X3 = clampInt(t.X3+int(rnd.NormFloat64()*16), -m, w-1+m)
			t.Y3 = clampInt(t.Y3+int(rnd.NormFloat64()*16), -m, h-1+m)
		}
		if t.Valid() {
			break
		}
	}
}

func (t *Triangle) Valid() bool {
	// Fast geometric validation without expensive trigonometry
	// Equivalent to checking that all angles are > 15 degrees
	
	x1, y1 := float64(t.X1), float64(t.Y1)
	x2, y2 := float64(t.X2), float64(t.Y2)
	x3, y3 := float64(t.X3), float64(t.Y3)
	
	// Calculate squared edge lengths to avoid sqrt
	dx12, dy12 := x2-x1, y2-y1
	dx23, dy23 := x3-x2, y3-y2
	dx31, dy31 := x1-x3, y1-y3
	
	len12Sq := dx12*dx12 + dy12*dy12
	len23Sq := dx23*dx23 + dy23*dy23
	len31Sq := dx31*dx31 + dy31*dy31
	
	// Check for degenerate triangle using area-based approach
	// Area = 0.5 * |cross product of two edges|
	// For non-degenerate triangle, area must be significant relative to perimeter
	crossProduct := dx12*dy31 - dy12*dx31
	areaSq := crossProduct * crossProduct
	
	// Minimum area threshold relative to edge lengths
	// This effectively filters out triangles with very small angles
	// Equivalent to checking angles > ~15 degrees without trigonometry
	minAreaSq := 0.07 * (len12Sq + len23Sq + len31Sq) // ~sin²(15°) ≈ 0.067
	
	if areaSq < minAreaSq {
		return false
	}
	
	// Additional check using dot products to ensure no angle is too small
	// For angle at vertex 1: cos(angle) = dot(v12, v13) / (|v12| * |v13|)
	// Angle > 15° means cos(angle) < cos(15°) ≈ 0.966
	// So cos²(angle) < 0.933, or dot² < 0.933 * len12Sq * len13Sq
	
	const maxCosSquared = 0.933 // cos²(15°)
	
	// Check angle at vertex 1
	dot12_13 := dx12*(-dx31) + dy12*(-dy31) // dot(v12, v13)
	if dot12_13*dot12_13 > maxCosSquared*len12Sq*len31Sq {
		return false
	}
	
	// Check angle at vertex 2  
	dot21_23 := (-dx12)*dx23 + (-dy12)*dy23 // dot(v21, v23)
	if dot21_23*dot21_23 > maxCosSquared*len12Sq*len23Sq {
		return false
	}
	
	// Check angle at vertex 3
	dot32_31 := (-dx23)*dx31 + (-dy23)*dy31 // dot(v32, v31)
	if dot32_31*dot32_31 > maxCosSquared*len23Sq*len31Sq {
		return false
	}
	
	return true
}

func (t *Triangle) Rasterize() []Scanline {
	buf := t.Worker.Lines[:0]
	lines := rasterizeTriangle(t.X1, t.Y1, t.X2, t.Y2, t.X3, t.Y3, buf)
	return cropScanlines(lines, t.Worker.W, t.Worker.H)
}

func rasterizeTriangle(x1, y1, x2, y2, x3, y3 int, buf []Scanline) []Scanline {
	if y1 > y3 {
		x1, x3 = x3, x1
		y1, y3 = y3, y1
	}
	if y1 > y2 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
	}
	if y2 > y3 {
		x2, x3 = x3, x2
		y2, y3 = y3, y2
	}
	if y2 == y3 {
		return rasterizeTriangleBottom(x1, y1, x2, y2, x3, y3, buf)
	} else if y1 == y2 {
		return rasterizeTriangleTop(x1, y1, x2, y2, x3, y3, buf)
	} else {
		x4 := x1 + int((float64(y2-y1)/float64(y3-y1))*float64(x3-x1))
		y4 := y2
		buf = rasterizeTriangleBottom(x1, y1, x2, y2, x4, y4, buf)
		buf = rasterizeTriangleTop(x2, y2, x4, y4, x3, y3, buf)
		return buf
	}
}

func rasterizeTriangleBottom(x1, y1, x2, y2, x3, y3 int, buf []Scanline) []Scanline {
	s1 := float64(x2-x1) / float64(y2-y1)
	s2 := float64(x3-x1) / float64(y3-y1)
	ax := float64(x1)
	bx := float64(x1)
	for y := y1; y <= y2; y++ {
		a := int(ax)
		b := int(bx)
		ax += s1
		bx += s2
		if a > b {
			a, b = b, a
		}
		buf = append(buf, Scanline{y, a, b, 0xffff})
	}
	return buf
}

func rasterizeTriangleTop(x1, y1, x2, y2, x3, y3 int, buf []Scanline) []Scanline {
	s1 := float64(x3-x1) / float64(y3-y1)
	s2 := float64(x3-x2) / float64(y3-y2)
	ax := float64(x3)
	bx := float64(x3)
	for y := y3; y > y1; y-- {
		ax -= s1
		bx -= s2
		a := int(ax)
		b := int(bx)
		if a > b {
			a, b = b, a
		}
		buf = append(buf, Scanline{y, a, b, 0xffff})
	}
	return buf
}
