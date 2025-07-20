package primitive

import (
	"fmt"
	"math"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/raster"
)

type Ellipse struct {
	Worker *Worker
	X, Y   int
	Rx, Ry int
	Circle bool
}

func NewRandomEllipse(worker *Worker) *Ellipse {
	rnd := worker.Rnd
	x := rnd.Intn(worker.W)
	y := rnd.Intn(worker.H)
	rx := rnd.Intn(32) + 1
	ry := rnd.Intn(32) + 1
	return &Ellipse{worker, x, y, rx, ry, false}
}

func NewRandomCircle(worker *Worker) *Ellipse {
	rnd := worker.Rnd
	x := rnd.Intn(worker.W)
	y := rnd.Intn(worker.H)
	r := rnd.Intn(32) + 1
	return &Ellipse{worker, x, y, r, r, true}
}

func (c *Ellipse) Draw(dc *gg.Context, scale float64) {
	dc.DrawEllipse(float64(c.X), float64(c.Y), float64(c.Rx), float64(c.Ry))
	dc.Fill()
}

func (c *Ellipse) SVG(attrs string) string {
	return fmt.Sprintf(
		"<ellipse %s cx=\"%d\" cy=\"%d\" rx=\"%d\" ry=\"%d\" />",
		attrs, c.X, c.Y, c.Rx, c.Ry)
}

func (c *Ellipse) Copy() Shape {
	a := *c
	return &a
}

func (c *Ellipse) Mutate() {
	w := c.Worker.W
	h := c.Worker.H
	rnd := c.Worker.Rnd
	
	// Pre-compute random values to avoid repeated function calls
	randChoice := rnd.Intn(3)
	
	switch randChoice {
	case 0:
		// Position mutation with bounds checking
		dx := int(rnd.NormFloat64() * 16)
		dy := int(rnd.NormFloat64() * 16)
		c.X = clampInt(c.X+dx, 0, w-1)
		c.Y = clampInt(c.Y+dy, 0, h-1)
	case 1:
		// Rx mutation
		dr := int(rnd.NormFloat64() * 16)
		c.Rx = clampInt(c.Rx+dr, 1, w-1)
		if c.Circle {
			c.Ry = c.Rx
		}
	case 2:
		// Ry mutation
		dr := int(rnd.NormFloat64() * 16)
		c.Ry = clampInt(c.Ry+dr, 1, h-1)
		if c.Circle {
			c.Rx = c.Ry
		}
	}
}

func (c *Ellipse) Rasterize() []Scanline {
	w := c.Worker.W
	h := c.Worker.H
	lines := c.Worker.Lines[:0]
	
	// Pre-allocate with estimated capacity to reduce allocations
	if cap(lines) < c.Ry*2 {
		lines = make([]Scanline, 0, c.Ry*2)
		c.Worker.Lines = lines
	}
	
	// Pre-compute constants to avoid repeated calculations
	rx := c.Rx
	ry := c.Ry
	rx2 := rx * rx
	ry2 := ry * ry
	aspect := float64(rx) / float64(ry)
	
	// Early bounds checking - skip if ellipse completely outside image
	if c.X+rx < 0 || c.X-rx >= w || c.Y+ry < 0 || c.Y-ry >= h {
		return lines
	}
	
	// Optimize loop bounds to only process visible scanlines
	maxDy := ry
	if c.Y-ry < 0 {
		// Adjust starting point if ellipse extends above image
		if c.Y+ry < 0 {
			return lines // completely above image
		}
	}
	if c.Y+ry >= h {
		// Adjust ending point if ellipse extends below image
		if c.Y-ry >= h {
			return lines // completely below image
		}
		maxDy = h - c.Y - 1
		if maxDy > ry {
			maxDy = ry
		}
	}
	
	for dy := 0; dy < maxDy; dy++ {
		y1 := c.Y - dy
		y2 := c.Y + dy
		
		// Use integer arithmetic where possible
		dySquared := dy * dy
		if dySquared >= ry2 {
			break // Outside ellipse bounds
		}
		
		// Optimized sqrt calculation using pre-computed values
		radius := int(math.Sqrt(float64(ry2-dySquared)) * aspect)
		x1 := c.X - radius
		x2 := c.X + radius
		
		// Inline bounds checking for better performance
		if x1 < 0 {
			x1 = 0
		}
		if x2 >= w {
			x2 = w - 1
		}
		
		// Skip invalid scanlines
		if x1 > x2 {
			continue
		}
		
		// Add scanlines with bounds checking
		if y1 >= 0 && y1 < h {
			lines = append(lines, Scanline{y1, x1, x2, 0xffff})
		}
		if y2 >= 0 && y2 < h && dy > 0 {
			lines = append(lines, Scanline{y2, x1, x2, 0xffff})
		}
	}
	
	c.Worker.Lines = lines
	return lines
}

type RotatedEllipse struct {
	Worker *Worker
	X, Y   float64
	Rx, Ry float64
	Angle  float64
}

func NewRandomRotatedEllipse(worker *Worker) *RotatedEllipse {
	rnd := worker.Rnd
	x := rnd.Float64() * float64(worker.W)
	y := rnd.Float64() * float64(worker.H)
	rx := rnd.Float64()*32 + 1
	ry := rnd.Float64()*32 + 1
	a := rnd.Float64() * 360
	return &RotatedEllipse{worker, x, y, rx, ry, a}
}

func (c *RotatedEllipse) Draw(dc *gg.Context, scale float64) {
	dc.Push()
	dc.RotateAbout(radians(c.Angle), c.X, c.Y)
	dc.DrawEllipse(c.X, c.Y, c.Rx, c.Ry)
	dc.Fill()
	dc.Pop()
}

func (c *RotatedEllipse) SVG(attrs string) string {
	return fmt.Sprintf(
		"<g transform=\"translate(%f %f) rotate(%f) scale(%f %f)\"><ellipse %s cx=\"0\" cy=\"0\" rx=\"1\" ry=\"1\" /></g>",
		c.X, c.Y, c.Angle, c.Rx, c.Ry, attrs)
}

func (c *RotatedEllipse) Copy() Shape {
	a := *c
	return &a
}

func (c *RotatedEllipse) Mutate() {
	w := c.Worker.W
	h := c.Worker.H
	rnd := c.Worker.Rnd
	
	// Pre-compute random values
	randChoice := rnd.Intn(3)
	
	switch randChoice {
	case 0:
		// Position mutation
		dx := rnd.NormFloat64() * 16
		dy := rnd.NormFloat64() * 16
		c.X = clamp(c.X+dx, 0, float64(w-1))
		c.Y = clamp(c.Y+dy, 0, float64(h-1))
	case 1:
		// Size mutation
		drx := rnd.NormFloat64() * 16
		dry := rnd.NormFloat64() * 16
		c.Rx = clamp(c.Rx+drx, 1, float64(w-1))
		c.Ry = clamp(c.Ry+dry, 1, float64(h-1)) // Fix: use h-1 for Ry
	case 2:
		// Angle mutation
		da := rnd.NormFloat64() * 32
		c.Angle = c.Angle + da
		// Normalize angle to prevent overflow
		for c.Angle > 360 {
			c.Angle -= 360
		}
		for c.Angle < 0 {
			c.Angle += 360
		}
	}
}

func (c *RotatedEllipse) Rasterize() []Scanline {
	// Early bounds checking for rotated ellipse
	maxRadius := math.Max(c.Rx, c.Ry)
	w := c.Worker.W
	h := c.Worker.H
	if c.X+maxRadius < 0 || c.X-maxRadius >= float64(w) || 
	   c.Y+maxRadius < 0 || c.Y-maxRadius >= float64(h) {
		return c.Worker.Lines[:0]
	}
	
	// Use fewer segments for small ellipses to improve performance
	n := 16
	if maxRadius < 16 {
		n = 8
	} else if maxRadius < 8 {
		n = 6
	}
	
	var path raster.Path
	
	// Pre-compute angle values to avoid repeated calculations
	angleRad := radians(c.Angle)
	sinAngle := math.Sin(angleRad)
	cosAngle := math.Cos(angleRad)
	
	// Pre-compute step size
	stepSize := 2 * math.Pi / float64(n)
	
	for i := 0; i < n; i++ {
		a1 := float64(i) * stepSize
		a2 := float64(i+1) * stepSize
		
		// Pre-compute trig values
		cos1, sin1 := math.Cos(a1), math.Sin(a1)
		cos2, sin2 := math.Cos(a2), math.Sin(a2)
		cosMid, sinMid := math.Cos(a1+(a2-a1)/2), math.Sin(a1+(a2-a1)/2)
		
		// Compute ellipse points
		x0 := c.Rx * cos1
		y0 := c.Ry * sin1
		x1 := c.Rx * cosMid
		y1 := c.Ry * sinMid
		x2 := c.Rx * cos2
		y2 := c.Ry * sin2
		
		// Bezier control point
		cx := 2*x1 - x0/2 - x2/2
		cy := 2*y1 - y0/2 - y2/2
		
		// Apply rotation using pre-computed sin/cos
		x0Rot := x0*cosAngle - y0*sinAngle
		y0Rot := x0*sinAngle + y0*cosAngle
		cxRot := cx*cosAngle - cy*sinAngle
		cyRot := cx*sinAngle + cy*cosAngle
		x2Rot := x2*cosAngle - y2*sinAngle
		y2Rot := x2*sinAngle + y2*cosAngle
		
		if i == 0 {
			path.Start(fixp(x0Rot+c.X, y0Rot+c.Y))
		}
		path.Add2(fixp(cxRot+c.X, cyRot+c.Y), fixp(x2Rot+c.X, y2Rot+c.Y))
	}
	return fillPath(c.Worker, path)
}
