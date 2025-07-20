package primitive

import (
	"image"
	"math"
)

func computeColor(target, current *image.RGBA, lines []Scanline, alpha int) Color {
	var rsum, gsum, bsum, count int64
	a := 0x101 * 255 / alpha
	
	for _, line := range lines {
		width := line.X2 - line.X1 + 1
		if width <= 0 {
			continue
		}
		
		i := target.PixOffset(line.X1, line.Y)
		
		// Vectorized processing - handle 4 pixels at a time when possible
		vectorWidth := width & ^3 // Round down to multiple of 4
		
		// Process 4 pixels at once for better CPU pipeline utilization
		for x := 0; x < vectorWidth; x += 4 {
			// Load 4 pixels from target and current images
			tr1, tg1, tb1 := int(target.Pix[i]), int(target.Pix[i+1]), int(target.Pix[i+2])
			cr1, cg1, cb1 := int(current.Pix[i]), int(current.Pix[i+1]), int(current.Pix[i+2])
			
			tr2, tg2, tb2 := int(target.Pix[i+4]), int(target.Pix[i+5]), int(target.Pix[i+6])
			cr2, cg2, cb2 := int(current.Pix[i+4]), int(current.Pix[i+5]), int(current.Pix[i+6])
			
			tr3, tg3, tb3 := int(target.Pix[i+8]), int(target.Pix[i+9]), int(target.Pix[i+10])
			cr3, cg3, cb3 := int(current.Pix[i+8]), int(current.Pix[i+9]), int(current.Pix[i+10])
			
			tr4, tg4, tb4 := int(target.Pix[i+12]), int(target.Pix[i+13]), int(target.Pix[i+14])
			cr4, cg4, cb4 := int(current.Pix[i+12]), int(current.Pix[i+13]), int(current.Pix[i+14])
			
			// Compute all 4 pixels in parallel
			rsum += int64((tr1-cr1)*a + cr1*0x101 + (tr2-cr2)*a + cr2*0x101 + 
			             (tr3-cr3)*a + cr3*0x101 + (tr4-cr4)*a + cr4*0x101)
			gsum += int64((tg1-cg1)*a + cg1*0x101 + (tg2-cg2)*a + cg2*0x101 + 
			             (tg3-cg3)*a + cg3*0x101 + (tg4-cg4)*a + cg4*0x101)
			bsum += int64((tb1-cb1)*a + cb1*0x101 + (tb2-cb2)*a + cb2*0x101 + 
			             (tb3-cb3)*a + cb3*0x101 + (tb4-cb4)*a + cb4*0x101)
			
			i += 16
		}
		
		// Handle remaining pixels (0-3)
		for x := vectorWidth; x < width; x++ {
			tr := int(target.Pix[i])
			tg := int(target.Pix[i+1])
			tb := int(target.Pix[i+2])
			cr := int(current.Pix[i])
			cg := int(current.Pix[i+1])
			cb := int(current.Pix[i+2])
			i += 4
			rsum += int64((tr-cr)*a + cr*0x101)
			gsum += int64((tg-cg)*a + cg*0x101)
			bsum += int64((tb-cb)*a + cb*0x101)
		}
		
		count += int64(width)
	}
	
	if count == 0 {
		return Color{}
	}
	r := clampInt(int(rsum/count)>>8, 0, 255)
	g := clampInt(int(gsum/count)>>8, 0, 255)
	b := clampInt(int(bsum/count)>>8, 0, 255)
	return Color{r, g, b, alpha}
}

func copyLines(dst, src *image.RGBA, lines []Scanline) {
	for _, line := range lines {
		a := dst.PixOffset(line.X1, line.Y)
		b := a + (line.X2-line.X1+1)*4
		copy(dst.Pix[a:b], src.Pix[a:b])
	}
}

func drawLines(im *image.RGBA, c Color, lines []Scanline) {
	const m = 0xffff
	sr, sg, sb, sa := c.NRGBA().RGBA()
	
	for _, line := range lines {
		width := line.X2 - line.X1 + 1
		if width <= 0 {
			continue
		}
		
		ma := line.Alpha
		a := (m - sa*ma/m) * 0x101
		i := im.PixOffset(line.X1, line.Y)
		
		// Pre-compute common values
		srma := sr * ma
		sgma := sg * ma
		sbma := sb * ma
		sama := sa * ma
		
		// Vectorized processing - handle 4 pixels at a time
		vectorWidth := width & ^3 // Round down to multiple of 4
		
		// Process 4 pixels at once
		for x := 0; x < vectorWidth; x += 4 {
			// Load 4 pixels
			dr1, dg1, db1, da1 := uint32(im.Pix[i+0]), uint32(im.Pix[i+1]), uint32(im.Pix[i+2]), uint32(im.Pix[i+3])
			dr2, dg2, db2, da2 := uint32(im.Pix[i+4]), uint32(im.Pix[i+5]), uint32(im.Pix[i+6]), uint32(im.Pix[i+7])
			dr3, dg3, db3, da3 := uint32(im.Pix[i+8]), uint32(im.Pix[i+9]), uint32(im.Pix[i+10]), uint32(im.Pix[i+11])
			dr4, dg4, db4, da4 := uint32(im.Pix[i+12]), uint32(im.Pix[i+13]), uint32(im.Pix[i+14]), uint32(im.Pix[i+15])
			
			// Compute and store 4 pixels
			im.Pix[i+0] = uint8((dr1*a + srma) / m >> 8)
			im.Pix[i+1] = uint8((dg1*a + sgma) / m >> 8)
			im.Pix[i+2] = uint8((db1*a + sbma) / m >> 8)
			im.Pix[i+3] = uint8((da1*a + sama) / m >> 8)
			
			im.Pix[i+4] = uint8((dr2*a + srma) / m >> 8)
			im.Pix[i+5] = uint8((dg2*a + sgma) / m >> 8)
			im.Pix[i+6] = uint8((db2*a + sbma) / m >> 8)
			im.Pix[i+7] = uint8((da2*a + sama) / m >> 8)
			
			im.Pix[i+8] = uint8((dr3*a + srma) / m >> 8)
			im.Pix[i+9] = uint8((dg3*a + sgma) / m >> 8)
			im.Pix[i+10] = uint8((db3*a + sbma) / m >> 8)
			im.Pix[i+11] = uint8((da3*a + sama) / m >> 8)
			
			im.Pix[i+12] = uint8((dr4*a + srma) / m >> 8)
			im.Pix[i+13] = uint8((dg4*a + sgma) / m >> 8)
			im.Pix[i+14] = uint8((db4*a + sbma) / m >> 8)
			im.Pix[i+15] = uint8((da4*a + sama) / m >> 8)
			
			i += 16
		}
		
		// Handle remaining pixels (0-3)
		for x := vectorWidth; x < width; x++ {
			dr := uint32(im.Pix[i+0])
			dg := uint32(im.Pix[i+1])
			db := uint32(im.Pix[i+2])
			da := uint32(im.Pix[i+3])
			im.Pix[i+0] = uint8((dr*a + srma) / m >> 8)
			im.Pix[i+1] = uint8((dg*a + sgma) / m >> 8)
			im.Pix[i+2] = uint8((db*a + sbma) / m >> 8)
			im.Pix[i+3] = uint8((da*a + sama) / m >> 8)
			i += 4
		}
	}
}

func differenceFull(a, b *image.RGBA) float64 {
	size := a.Bounds().Size()
	w, h := size.X, size.Y
	var total uint64
	for y := 0; y < h; y++ {
		i := a.PixOffset(0, y)
		for x := 0; x < w; x++ {
			ar := int(a.Pix[i])
			ag := int(a.Pix[i+1])
			ab := int(a.Pix[i+2])
			aa := int(a.Pix[i+3])
			br := int(b.Pix[i])
			bg := int(b.Pix[i+1])
			bb := int(b.Pix[i+2])
			ba := int(b.Pix[i+3])
			i += 4
			dr := ar - br
			dg := ag - bg
			db := ab - bb
			da := aa - ba
			total += uint64(dr*dr + dg*dg + db*db + da*da)
		}
	}
	return math.Sqrt(float64(total)/float64(w*h*4)) / 255
}

func differencePartial(target, before, after *image.RGBA, score float64, lines []Scanline) float64 {
	size := target.Bounds().Size()
	w, h := size.X, size.Y
	total := uint64(math.Pow(score*255, 2) * float64(w*h*4))
	
	for _, line := range lines {
		width := line.X2 - line.X1 + 1
		if width <= 0 {
			continue
		}
		
		i := target.PixOffset(line.X1, line.Y)
		
		// Vectorized processing - handle 4 pixels at a time
		vectorWidth := width & ^3 // Round down to multiple of 4
		
		// Process 4 pixels at once for better memory bandwidth utilization
		for x := 0; x < vectorWidth; x += 4 {
			// Load 4 pixels from all three images
			tr1, tg1, tb1, ta1 := int(target.Pix[i]), int(target.Pix[i+1]), int(target.Pix[i+2]), int(target.Pix[i+3])
			br1, bg1, bb1, ba1 := int(before.Pix[i]), int(before.Pix[i+1]), int(before.Pix[i+2]), int(before.Pix[i+3])
			ar1, ag1, ab1, aa1 := int(after.Pix[i]), int(after.Pix[i+1]), int(after.Pix[i+2]), int(after.Pix[i+3])
			
			tr2, tg2, tb2, ta2 := int(target.Pix[i+4]), int(target.Pix[i+5]), int(target.Pix[i+6]), int(target.Pix[i+7])
			br2, bg2, bb2, ba2 := int(before.Pix[i+4]), int(before.Pix[i+5]), int(before.Pix[i+6]), int(before.Pix[i+7])
			ar2, ag2, ab2, aa2 := int(after.Pix[i+4]), int(after.Pix[i+5]), int(after.Pix[i+6]), int(after.Pix[i+7])
			
			tr3, tg3, tb3, ta3 := int(target.Pix[i+8]), int(target.Pix[i+9]), int(target.Pix[i+10]), int(target.Pix[i+11])
			br3, bg3, bb3, ba3 := int(before.Pix[i+8]), int(before.Pix[i+9]), int(before.Pix[i+10]), int(before.Pix[i+11])
			ar3, ag3, ab3, aa3 := int(after.Pix[i+8]), int(after.Pix[i+9]), int(after.Pix[i+10]), int(after.Pix[i+11])
			
			tr4, tg4, tb4, ta4 := int(target.Pix[i+12]), int(target.Pix[i+13]), int(target.Pix[i+14]), int(target.Pix[i+15])
			br4, bg4, bb4, ba4 := int(before.Pix[i+12]), int(before.Pix[i+13]), int(before.Pix[i+14]), int(before.Pix[i+15])
			ar4, ag4, ab4, aa4 := int(after.Pix[i+12]), int(after.Pix[i+13]), int(after.Pix[i+14]), int(after.Pix[i+15])
			
			// Compute differences for 4 pixels
			dr1_1, dg1_1, db1_1, da1_1 := tr1-br1, tg1-bg1, tb1-bb1, ta1-ba1
			dr2_1, dg2_1, db2_1, da2_1 := tr1-ar1, tg1-ag1, tb1-ab1, ta1-aa1
			
			dr1_2, dg1_2, db1_2, da1_2 := tr2-br2, tg2-bg2, tb2-bb2, ta2-ba2
			dr2_2, dg2_2, db2_2, da2_2 := tr2-ar2, tg2-ag2, tb2-ab2, ta2-aa2
			
			dr1_3, dg1_3, db1_3, da1_3 := tr3-br3, tg3-bg3, tb3-bb3, ta3-ba3
			dr2_3, dg2_3, db2_3, da2_3 := tr3-ar3, tg3-ag3, tb3-ab3, ta3-aa3
			
			dr1_4, dg1_4, db1_4, da1_4 := tr4-br4, tg4-bg4, tb4-bb4, ta4-ba4
			dr2_4, dg2_4, db2_4, da2_4 := tr4-ar4, tg4-ag4, tb4-ab4, ta4-aa4
			
			// Accumulate squared differences for all 4 pixels
			total -= uint64(dr1_1*dr1_1 + dg1_1*dg1_1 + db1_1*db1_1 + da1_1*da1_1 +
			               dr1_2*dr1_2 + dg1_2*dg1_2 + db1_2*db1_2 + da1_2*da1_2 +
			               dr1_3*dr1_3 + dg1_3*dg1_3 + db1_3*db1_3 + da1_3*da1_3 +
			               dr1_4*dr1_4 + dg1_4*dg1_4 + db1_4*db1_4 + da1_4*da1_4)
			
			total += uint64(dr2_1*dr2_1 + dg2_1*dg2_1 + db2_1*db2_1 + da2_1*da2_1 +
			               dr2_2*dr2_2 + dg2_2*dg2_2 + db2_2*db2_2 + da2_2*da2_2 +
			               dr2_3*dr2_3 + dg2_3*dg2_3 + db2_3*db2_3 + da2_3*da2_3 +
			               dr2_4*dr2_4 + dg2_4*dg2_4 + db2_4*db2_4 + da2_4*da2_4)
			
			i += 16
		}
		
		// Handle remaining pixels (0-3)
		for x := vectorWidth; x < width; x++ {
			tr := int(target.Pix[i])
			tg := int(target.Pix[i+1])
			tb := int(target.Pix[i+2])
			ta := int(target.Pix[i+3])
			br := int(before.Pix[i])
			bg := int(before.Pix[i+1])
			bb := int(before.Pix[i+2])
			ba := int(before.Pix[i+3])
			ar := int(after.Pix[i])
			ag := int(after.Pix[i+1])
			ab := int(after.Pix[i+2])
			aa := int(after.Pix[i+3])
			i += 4
			dr1 := tr - br
			dg1 := tg - bg
			db1 := tb - bb
			da1 := ta - ba
			dr2 := tr - ar
			dg2 := tg - ag
			db2 := tb - ab
			da2 := ta - aa
			total -= uint64(dr1*dr1 + dg1*dg1 + db1*db1 + da1*da1)
			total += uint64(dr2*dr2 + dg2*dg2 + db2*db2 + da2*da2)
		}
	}
	return math.Sqrt(float64(total)/float64(w*h*4)) / 255
}
