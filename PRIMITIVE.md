# PRIMITIVE.md

## Go File Analysis and Performance Optimization Guide

This document provides a comprehensive analysis of each Go file in the `primitive/` folder, describing their functionality and suggesting performance improvements for the primitive image generation algorithm.

---

## color.go

**Purpose**: Defines custom Color struct and handles color conversion between Go's standard color types and the primitive algorithm's integer-based RGBA format.

**Key Functions**:
- `Color` struct: Efficient integer-based RGBA representation (0-255 range)
- `MakeColor()`: Converts standard Go color types to primitive Color format
- `MakeHexColor()`: Parses hex color strings (#RGB, #RRGGBB, etc.)
- `NRGBA()`: Converts back to Go's standard color format

**Performance Optimizations**:
1. **Replace `fmt.Sscanf` with manual hex parsing** - Eliminate string formatting overhead
2. **Implement color pools** - Reduce GC pressure for frequently used colors
3. **Cache `NRGBA()` results** - Precompute for common colors like backgrounds
4. **Add SIMD support** - Batch color operations for multiple colors simultaneously

---

## core.go

**Purpose**: Core image processing functions for pixel-level operations, including color computation, image difference calculations, and scanline-based drawing.

**Key Functions**:
- `computeColor()`: Optimal color calculation for shapes using alpha blending
- `drawLines()`: Alpha-blended scanline rendering with fixed-point arithmetic
- `differenceFull()`: Complete image difference calculation (RMS)
- `differencePartial()`: Incremental difference updates for optimization

**Performance Optimizations**:
1. **SIMD vectorization** - Parallel processing of RGBA channels
2. **Memory access optimization** - Pre-calculate `PixOffset` values, use pointer arithmetic
3. **Parallel processing** - Split scanlines across goroutines for large shapes
4. **Lookup tables** - Pre-compute alpha blending coefficients
5. **Integer arithmetic** - Replace floating-point operations where possible

---

## ellipse.go

**Purpose**: Implements ellipse and rotated ellipse shapes with random generation, mutation, validation, and rasterization capabilities.

**Key Functions**:
- `NewRandomEllipse()`: Creates random ellipses with constrained dimensions
- `Rasterize()`: Mathematical ellipse-to-scanline conversion
- `RotatedEllipse.Rasterize()`: Uses 16-segment Bézier approximation

**Performance Optimizations**:
1. **Pre-compute constants** - Cache `Ry*Ry` and aspect ratios outside loops
2. **Reduce Bézier segments** - Use fewer segments for small ellipses
3. **Object pooling** - Reuse ellipse instances and scanline buffers
4. **SIMD scanline generation** - Vectorize ellipse equation calculations

---

## heatmap.go

**Purpose**: Tracks pixel frequency during shape rendering to optimize placement and ordering of primitive shapes.

**Key Functions**:
- `Add()`: Accumulates pixel coverage from scanlines
- `Image()`: Converts heatmap to visual grayscale representation with gamma correction
- `AddHeatmap()`: Merges multiple heatmaps

**Performance Optimizations**:
1. **Vectorized operations** - Use SIMD for bulk array operations
2. **Gamma lookup tables** - Pre-compute gamma correction values
3. **Parallel processing** - Split image conversion across goroutines
4. **Memory zeroing optimization** - Use unsafe operations for faster clearing

---

## log.go

**Purpose**: Simple level-based logging system with hierarchical indentation for tracking optimization progress.

**Key Functions**:
- `Log()`: Core logging with level filtering
- `v()`, `vv()`, `vvv()`: Convenience functions with automatic indentation

**Performance Optimizations**:
1. **Fast path for disabled logging** - Early return when LogLevel is insufficient
2. **Pre-compute format strings** - Avoid string concatenation in hot paths
3. **Avoid Printf overhead** - Use Print for simple strings without formatting
4. **Compile-time optimization** - Use build tags to eliminate logging in production

---

## model.go

**Purpose**: Core Model struct that orchestrates the primitive generation algorithm, managing state and coordinating parallel workers.

**Key Functions**:
- `NewModel()`: Initialize model with target image and worker pool
- `Step()`: Single optimization iteration with parallel worker coordination
- `Add()`: Adds optimized shape to current image
- `runWorkers()`: Parallel shape optimization across multiple goroutines

**Performance Optimizations**:
1. **Memory pools** - Reuse image buffers and temporary objects
2. **Worker persistence** - Maintain worker pools across iterations
3. **Batch operations** - Group multiple small shapes
4. **Data locality improvements** - Reorganize struct fields for cache efficiency
5. **Atomic operations** - Replace channels with atomic counters where possible

---

## optimize.go

**Purpose**: Implements optimization algorithms (hill climbing, simulated annealing) for shape parameter tuning.

**Key Functions**:
- `HillClimb()`: Deterministic optimization accepting only improvements
- `Anneal()`: Simulated annealing with temperature-based acceptance
- `PreAnneal()`: Temperature calibration through random sampling

**Performance Optimizations**:
1. **Reduce copy operations** - Implement copy-on-write semantics
2. **Adaptive algorithms** - Dynamic mutation rates and temperature schedules
3. **Parallel optimization chains** - Multiple concurrent optimization runs
4. **Generics usage** - Eliminate interface overhead with Go 1.18+ generics

---

## polygon.go

**Purpose**: Variable-vertex polygon implementation supporting both convex and non-convex shapes.

**Key Functions**:
- `NewRandomPolygon()`: Creates random polygons with specified vertex count
- `Mutate()`: Vertex swapping or position adjustment with validation
- `Valid()`: Convexity checking using cross products
- `Rasterize()`: Polygon-to-scanline conversion

**Performance Optimizations**:
1. **String builder for SVG** - Eliminate string concatenation allocations
2. **Coordinate array pooling** - Reuse float64 slices via sync.Pool
3. **Cached validity checks** - Store validation state
4. **Iteration limits** - Prevent infinite loops in constrained mutations

---

## quadratic.go

**Purpose**: Quadratic Bézier curve implementation as primitive shapes for approximating complex curves.

**Key Functions**:
- `NewRandomQuadratic()`: Random curve generation with control points
- `Mutate()`: Control point and width mutations with validation
- `Valid()`: Geometric validation ensuring meaningful curve shape
- `Rasterize()`: Curve-to-scanline conversion using freetype raster

**Performance Optimizations**:
1. **Fix mutation bug** - Change `rnd.Intn(3)` to `rnd.Intn(4)` to include width mutations
2. **Pre-allocate paths** - Reuse raster.Path objects
3. **Bounds-checked mutations** - Reduce invalid mutation attempts
4. **Distance calculation caching** - Cache validation computations

---

## raster.go

**Purpose**: Bridge between vector graphics and scanline representation, converting paths to pixel spans.

**Key Functions**:
- `fillPath()`: Filled shape rasterization using non-zero winding rule
- `strokePath()`: Stroked path rasterization with width and line caps
- `fix()`: Fixed-point coordinate conversion for precise rendering

**Performance Optimizations**:
1. **Scanline slice pooling** - Implement sync.Pool for scanline reuse
2. **Batch rasterization** - Process multiple similar shapes together
3. **Bounds checking optimization** - Early termination for out-of-bounds spans
4. **Capacity estimation** - Dynamic sizing based on shape complexity

---

## rectangle.go

**Purpose**: Implements axis-aligned and rotated rectangle shapes with efficient rasterization.

**Key Functions**:
- `Rectangle`/`RotatedRectangle`: Two rectangle variants
- `Rasterize()`: Direct scanline generation for axis-aligned rectangles
- `RotatedRectangle.Rasterize()`: Complex edge tracing for rotated rectangles

**Performance Optimizations**:
1. **Buffer reuse** - Reuse temporary arrays in rotated rectangle rasterization
2. **Bresenham's algorithm** - Replace brute-force edge sampling
3. **Trigonometric caching** - Cache sin/cos for common rotation angles
4. **Integer arithmetic** - Use fixed-point math where precision allows

---

## scanline.go

**Purpose**: Core scanline data structure and clipping operations for efficient pixel-level rendering.

**Key Functions**:
- `Scanline` struct: Horizontal line segment with Y, X1, X2, Alpha
- `cropScanlines()`: Clips scanlines to image boundaries

**Performance Optimizations**:
1. **Index-based iteration** - Avoid struct copying in cropScanlines
2. **SIMD-friendly data layout** - Structure-of-arrays for batch operations
3. **Branch prediction** - Reorder conditions by likelihood
4. **Inline bounds checking** - Eliminate function call overhead

---

## shape.go

**Purpose**: Defines the Shape interface and type enumeration for polymorphic shape handling.

**Key Functions**:
- `Shape` interface: Common operations (Rasterize, Copy, Mutate, Draw, SVG)
- `ShapeType` enum: 8 different primitive shape types

**Performance Optimizations**:
1. **Scanline buffer pooling** - Shared buffer management across shape types
2. **SIMD rasterization** - Vectorized operations for shape-specific algorithms
3. **Generics implementation** - Type-safe optimization avoiding interface overhead
4. **Shape-specific optimizations** - Lookup tables for small shapes, approximations for large ones

---

## state.go

**Purpose**: Represents optimization state implementing the Annealable interface for hill climbing and simulated annealing.

**Key Functions**:
- `Energy()`: Lazy fitness evaluation with memoization
- `DoMove()`/`UndoMove()`: Mutation and rollback operations
- `Copy()`: Deep state copying for algorithm branches

**Performance Optimizations**:
1. **Object pooling** - Reuse State structs to reduce GC pressure
2. **Smart alpha mutation** - Adaptive step sizes instead of fixed ±10 range
3. **Batch operations** - Group multiple mutations for amortized cost
4. **Safe type handling** - Add bounds checking for undo operations

---

## triangle.go

**Purpose**: Triangle shape implementation with three-vertex geometry and scanline rasterization.

**Key Functions**:
- `NewRandomTriangle()`: Random triangle generation with proximity constraints
- `Valid()`: Angle-based validation preventing degenerate triangles
- `Rasterize()`: Triangle decomposition into top-flat and bottom-flat cases

**Performance Optimizations**:
1. **Replace angle validation** - Use area-based or cross-product validation
2. **Pre-allocate scanline capacity** - Estimate based on triangle size
3. **Integer-only rasterization** - Avoid floating-point arithmetic in hot paths
4. **SIMD edge walking** - Vectorize slope calculations and scanline generation

---

## util.go

**Purpose**: Utility functions for I/O, image processing, mathematical operations, and format conversions.

**Key Functions**:
- `LoadImage()`/`SavePNG()`/`SaveJPG()`: Image I/O operations
- `AverageImageColor()`: Mean color calculation for background initialization
- `imageToRGBA()`: Format conversion for consistent processing
- `rotate()`: 2D coordinate transformation

**Performance Optimizations**:
1. **Fix AverageImageColor()** - Critical O(n²) to O(n) optimization using stride order
2. **RGBA conversion caching** - Skip conversion if already RGBA format
3. **Trigonometric lookup tables** - Cache sin/cos for rotation operations
4. **Memory pool for images** - Reuse RGBA images via sync.Pool

---

## worker.go

**Purpose**: Core computation engine managing shape generation, evaluation, and optimization for individual worker threads.

**Key Functions**:
- `Energy()`: Shape fitness evaluation through rasterization and difference calculation
- `BestHillClimbState()`: Combines random generation with hill climbing
- `BestRandomState()`: Parallel random state generation and selection

**Performance Optimizations**:
1. **Dynamic buffer sizing** - Base scanline capacity on image height
2. **Parallel state generation** - Use goroutines for embarrassingly parallel loops
3. **Buffer pooling** - Reuse image buffers to reduce GC pressure
4. **Early termination** - Stop optimization when improvement falls below threshold
5. **Memory layout optimization** - Reorganize Worker struct for cache efficiency

---

## Priority Optimization Recommendations

**Highest Impact (Critical Path)**:
1. **util.go: AverageImageColor()** - Fix O(n²) algorithm
2. **core.go: SIMD vectorization** - Parallel RGBA processing
3. **worker.go: Buffer management** - Reduce GC pressure in hot path
4. **triangle.go: Validation optimization** - Replace expensive trigonometry

**High Impact (Frequent Operations)**:
5. **scanline.go: Index-based iteration** - Eliminate struct copying
6. **model.go: Memory pools** - Reuse temporary objects
7. **All shape files: Object pooling** - Reduce allocation overhead

**Medium Impact (Algorithm Improvements)**:
8. **optimize.go: Parallel chains** - Multiple concurrent optimizations
9. **raster.go: Batch processing** - Group similar rasterization operations
10. **log.go: Compile-time elimination** - Remove logging overhead in production