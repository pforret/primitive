# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Primitive is a Go application that reproduces images using geometric primitives. It implements a hill-climbing optimization algorithm to iteratively add shapes (triangles, rectangles, ellipses, etc.) to approximate a target image while minimizing root-mean-square error.

## Build and Run Commands

This is a standard Go project. Use these commands:

```bash
# Build the application
go build -o primitive

# Run with basic usage
./primitive -i input.png -o output.png -n 100

# Install as a command (if GOPATH is set up)
go install
```

## Architecture

### Core Components

- **main.go**: CLI interface with flag parsing and main execution loop
- **primitive/**: Core package containing the algorithm implementation
  - **model.go**: Main Model struct that manages the optimization process
  - **shape.go**: Shape interface and ShapeType definitions  
  - **core.go**: Low-level image processing functions (color computation, scanline operations)
  - **worker.go**: Parallel worker implementation for shape optimization
  - **optimize.go**: Hill climbing and optimization algorithms

### Shape System

The application supports multiple primitive shapes through a common interface:

```go
type Shape interface {
    Rasterize() []Scanline
    Copy() Shape
    Mutate()
    Draw(dc *gg.Context, scale float64)
    SVG(attrs string) string
}
```

Shape types: Triangle, Rectangle, Ellipse, Circle, RotatedRectangle, Quadratic (Bezier), RotatedEllipse, Polygon

### Algorithm Flow

1. **Initialization**: Create Model with target image and background color
2. **Shape Generation**: Workers generate random shapes of specified type
3. **Optimization**: Hill climbing to find optimal shape parameters
4. **Evaluation**: Score shapes using RMSE against target image
5. **Integration**: Add best shape to current image and repeat

### Key Files

- **primitive/model.go:25**: Model creation and initialization
- **primitive/shape.go:5**: Shape interface definition
- **primitive/worker.go**: Parallel optimization workers
- **main.go:91**: Main execution loop and CLI argument handling

## Dependencies

- `github.com/fogleman/gg`: 2D graphics library for rendering
- `github.com/nfnt/resize`: Image resizing utilities

## Output Formats

- PNG/JPG: Raster output
- SVG: Vector output  
- GIF: Animated sequence (requires ImageMagick)