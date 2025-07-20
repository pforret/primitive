#!/usr/bin/env bash

# Primitive Performance Benchmark Script
# Compiles the Go binary and runs comprehensive performance tests
# Outputs results in format compatible with baseline comparison

set -e

# Configuration
SHAPES=500
BINARY_NAME="primitive-benchmark"
OUTPUT_DIR="/tmp/primitive_bench_output"
TIMESTAMP=$(date +"%Y-%m-%d %H:%M:%S")

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Primitive Performance Benchmark ===${NC}"
echo "Date: $TIMESTAMP"
echo "Shapes per test: $SHAPES"
echo "Go version: $(go version)"
echo "Platform: $(uname -s) $(uname -r)"
echo

# Step 1: Compile the binary
echo -e "${YELLOW}[1/4] Compiling Go binary...${NC}"
if go build -o "$BINARY_NAME"; then
    echo -e "${GREEN}✓ Binary compiled successfully: $BINARY_NAME${NC}"
else
    echo -e "${RED}✗ Compilation failed${NC}"
    exit 1
fi

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Step 2: Define test configurations
echo -e "${YELLOW}[2/4] Preparing test configurations...${NC}"

# Test images from examples/
IMAGES=(
    "lenna.png"
    "monalisa.png" 
    "owl.png"
    "pyramids.png"
)

echo -e "${GREEN}✓ Found ${#IMAGES[@]} test images, 3 shape modes${NC}"

# Function to get mode number
get_mode_num() {
    case $1 in
        "Triangle") echo "1" ;;
        "Rectangle") echo "2" ;;
        "Ellipse") echo "3" ;;
        "Circle") echo "4" ;;
        "Beziers") echo "6" ;;
    esac
}

# Step 3: Run benchmarks
echo -e "${YELLOW}[3/4] Running performance benchmarks...${NC}"
echo

# Results file
RESULTS_FILE="/tmp/benchmark_results.txt"
echo "" > "$RESULTS_FILE"

echo "## Benchmark Results - $SHAPES Shapes"
echo
echo "### Performance Summary (Wall-clock time in seconds)"
echo
echo "| Image        | Triangle (m=1) | Rectangle (m=2)| Ellipse (m=3)  | Circle (m=4)   | Beziers (m=6)  |"
echo "|--------------|----------------|----------------|----------------|----------------|----------------|"

for image in "${IMAGES[@]}"; do
    if [ ! -f "examples/$image" ]; then
        echo -e "${RED}Warning: examples/$image not found, skipping${NC}"
        continue
    fi
    
    printf "| %-12s |" "${image}"
    
    for mode_name in "Triangle" "Rectangle" "Ellipse" "Circle" "Beziers"; do
        mode_num=$(get_mode_num "$mode_name")
        output_file="$OUTPUT_DIR/${image%.*}_${mode_name,,}_$SHAPES.png"
        
        #echo -e "${BLUE}Testing: $image with $mode_name mode ($SHAPES shapes)${NC}"
        
        # Run benchmark with time measurement
        start_time=$(date +%s.%N)
        if ./"$BINARY_NAME" -i "examples/$image" -o "$output_file" -n "$SHAPES" -m "$mode_num" >/dev/null 2>&1; then
            end_time=$(date +%s.%N)
            wall_time=$(echo "$end_time - $start_time" | bc -l | awk '{printf "%.3f", $1}')
            
            # Store results
            echo "$image,$mode_name,$wall_time" >> "$RESULTS_FILE"
            
            printf " %-14s |" "${wall_time}s"
            # echo -e "${GREEN}  ✓ ${wall_time}s${NC}"
        else
            printf " %-14s |" "-"
            echo -e "${RED}  ✗ Failed${NC}"
            echo "$image,$mode_name,-" >> "$RESULTS_FILE"
        fi
    done
    echo
done

# Step 4: Generate comparison report
echo -e "${YELLOW}[4/4] Generating comparison report...${NC}"
echo
echo "### Benchmark Completion Summary"
echo

# Calculate averages for each mode
triangle_total=0
triangle_count=0
rectangle_total=0
rectangle_count=0
ellipse_total=0
ellipse_count=0

while IFS=',' read -r img mode time; do
    if [ "$time" != "-" ] && [ -n "$time" ]; then
        case $mode in
            "Triangle")
                triangle_total=$(echo "$triangle_total + $time" | bc -l)
                triangle_count=$((triangle_count + 1))
                ;;
            "Rectangle")
                rectangle_total=$(echo "$rectangle_total + $time" | bc -l)
                rectangle_count=$((rectangle_count + 1))
                ;;
            "Ellipse")
                ellipse_total=$(echo "$ellipse_total + $time" | bc -l)
                ellipse_count=$((ellipse_count + 1))
                ;;
        esac
    fi
done < "$RESULTS_FILE"

# Calculate and display averages
if [ $triangle_count -gt 0 ]; then
    triangle_avg=$(echo "scale=3; $triangle_total / $triangle_count" | bc -l)
    echo -e "${GREEN}Triangle mode average: ${triangle_avg}s (across $triangle_count images)${NC}"
fi

if [ $rectangle_count -gt 0 ]; then
    rectangle_avg=$(echo "scale=3; $rectangle_total / $rectangle_count" | bc -l)
    echo -e "${GREEN}Rectangle mode average: ${rectangle_avg}s (across $rectangle_count images)${NC}"
fi

if [ $ellipse_count -gt 0 ]; then
    ellipse_avg=$(echo "scale=3; $ellipse_total / $ellipse_count" | bc -l)
    echo -e "${GREEN}Ellipse mode average: ${ellipse_avg}s (across $ellipse_count images)${NC}"
fi

echo
echo "### Performance Hierarchy"
echo
echo "**Current ranking** (average time per 500 shapes):"

# Simple ranking based on calculated averages
echo "1. **Rectangle mode**: ${rectangle_avg:-N/A}s average"
echo "2. **Triangle mode**: ${triangle_avg:-N/A}s average"  
echo "3. **Ellipse mode**: ${ellipse_avg:-N/A}s average"

echo
echo "### Comparison with Baseline"
echo
echo "| Mode      | Baseline Avg | Current Avg | Difference | Change % |"
echo "|-----------|--------------|-------------|------------|----------|"

# Baseline values from BASELINE_500_RESULTS.md
baseline_triangle="8.201"
baseline_rectangle="6.149"
baseline_ellipse="13.829"

# Compare Triangle
if [ -n "$triangle_avg" ]; then
    diff=$(echo "scale=3; $triangle_avg - $baseline_triangle" | bc -l)
    pct=$(echo "scale=1; ($triangle_avg - $baseline_triangle) * 100 / $baseline_triangle" | bc -l)
    printf "| %-9s | %-12s | %-11s | %-10s | %s |\n" "Triangle" "${baseline_triangle}s" "${triangle_avg}s" "${diff}s" "${pct}%"
fi

# Compare Rectangle
if [ -n "$rectangle_avg" ]; then
    diff=$(echo "scale=3; $rectangle_avg - $baseline_rectangle" | bc -l)
    pct=$(echo "scale=1; ($rectangle_avg - $baseline_rectangle) * 100 / $baseline_rectangle" | bc -l)
    printf "| %-9s | %-12s | %-11s | %-10s | %s |\n" "Rectangle" "${baseline_rectangle}s" "${rectangle_avg}s" "${diff}s" "${pct}%"
fi

# Compare Ellipse
if [ -n "$ellipse_avg" ]; then
    diff=$(echo "scale=3; $ellipse_avg - $baseline_ellipse" | bc -l)
    pct=$(echo "scale=1; ($ellipse_avg - $baseline_ellipse) * 100 / $baseline_ellipse" | bc -l)
    printf "| %-9s | %-12s | %-11s | %-10s | %s |\n" "Ellipse" "${baseline_ellipse}s" "${ellipse_avg}s" "${diff}s" "${pct}%"
fi

echo
echo -e "${BLUE}=== Benchmark Complete ===${NC}"
echo "Output images saved to: $OUTPUT_DIR"
echo "Binary: $BINARY_NAME"
echo "Results: $RESULTS_FILE"
echo

# Cleanup option
read -p "Remove benchmark binary? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    rm -f "$BINARY_NAME"
    echo -e "${GREEN}✓ Binary removed${NC}"
fi

# Clean up results file
rm -f "$RESULTS_FILE"