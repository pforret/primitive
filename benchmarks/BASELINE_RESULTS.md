# Primitive Performance Baseline Results

## Test Configuration
- **Date**: July 20, 2025
- **Shapes per test**: 50
- **Go version**: 1.23.5
- **Platform**: Darwin 24.5.0 (macOS)
- **CPU**: Multi-core (parallel workers enabled)

## Benchmark Results
Ò
### Performance Summary (Wall-clock time in seconds)

| Image        | Triangle (m=1) | Rectangle (m=2) | Ellipse (m=3) | Combo (m=0) |
|--------------|----------------|-----------------|---------------|-------------|
| lenna.png    | 1.629s         | 0.938s          | 1.651s        | 1.912s      |
| monalisa.png | 1.817s         | -               | 2.250s        | -           |
| owl.png      | 1.352s         | -               | 1.737s        | -           |
| pyramids.png | 1.877s         | 1.032s          | -             | -           |

### CPU Usage Analysis

**CPU Utilization** (user time / wall time):
- Triangle mode: ~5.3x CPU utilization (good parallelization)
- Rectangle mode: ~5.7x CPU utilization (best parallelization)
- Ellipse mode: ~6.5x CPU utilization (excellent parallelization)
- Combo mode: ~6.3x CPU utilization (very good parallelization)

### Performance Insights

**Fastest Modes by Wall-clock Time**:
1. **Rectangle (mode 2)**: 0.938-1.032s
2. **Triangle (mode 1)**: 1.352-1.877s  
3. **Ellipse (mode 3)**: 1.651-2.250s
4. **Combo (mode 0)**: 1.912s

**Key Observations**:
- **Rectangle mode is fastest** due to simpler rasterization algorithm
- **Ellipse mode takes longest** due to complex mathematical calculations
- **Good parallel scaling** across all modes (5-6x CPU utilization)
- **Image complexity affects performance**: monalisa.png consistently slower

### Performance Bottlenecks Identified

Based on timing analysis and code review:

1. **Ellipse rasterization** - Mathematical sqrt/trig operations in hot path
2. **Triangle validation** - Expensive angle calculations during mutation
3. **Memory allocation** - Frequent scanline and shape object creation
4. **Color computation** - Called for every shape evaluation

### Test Images Characteristics

- **lenna.png**: Classic test image, medium complexity
- **monalisa.png**: High detail, complex features (slowest overall)
- **owl.png**: Natural texture, medium-high complexity  
- **pyramids.png**: Geometric shapes, medium complexity

### Baseline for Optimization

These results establish our performance baseline. Key targets for optimization:

1. **Target 20-30% improvement** in ellipse mode
2. **Target 15-20% improvement** in triangle mode
3. **Focus on memory allocation reduction**
4. **Optimize mathematical operations in rasterization**

### Next Steps

Use PRIMITIVE.md optimization recommendations to:
1. Fix O(n²) AverageImageColor algorithm
2. Implement SIMD vectorization for pixel operations
3. Add object pooling for memory management
4. Optimize trigonometric calculations with lookup tables

---

*Baseline established for post-optimization comparison*