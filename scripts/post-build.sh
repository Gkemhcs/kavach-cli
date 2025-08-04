#!/bin/bash

# Post-build script for GoReleaser
# This script runs after building the binaries

set -e

echo "🎉 Post-build process started..."

# List all built binaries
echo "📦 Built binaries:"
find dist -name "kavach*" -type f | sort

# Test each binary
echo "🧪 Testing binaries..."
for binary in $(find dist -name "kavach*" -type f); do
    echo "Testing: $binary"
    
    # Check if binary exists and is executable
    if [ ! -f "$binary" ]; then
        echo "❌ $binary - file does not exist"
        exit 1
    fi
    
    if [ ! -x "$binary" ]; then
        echo "❌ $binary - file is not executable"
        exit 1
    fi
    
    # Show binary info
    echo "📋 Binary info:"
    ls -la "$binary"
    file "$binary"
    
    # Try to run the binary with strace to see what's happening
    echo "🔍 Testing binary execution..."
    
    # First, try a simple execution
    if timeout 10s "$binary" > /dev/null 2>&1; then
        echo "✅ $binary - basic execution works"
    else
        echo "❌ $binary - basic execution failed"
        echo "🔍 Trying with strace to debug..."
        if command -v strace >/dev/null 2>&1; then
            timeout 10s strace -f -e trace=file,process "$binary" 2>&1 | head -20
        fi
    fi
    
    # Test version command
    echo "🔍 Testing version command..."
    if timeout 10s "$binary" version --short 2>&1; then
        echo "✅ $binary - version command works"
    else
        echo "❌ $binary - version command failed"
        echo "🔍 Trying without --short flag..."
        if timeout 10s "$binary" version 2>&1; then
            echo "✅ $binary - version command works without --short"
        else
            echo "❌ $binary - version command failed completely"
            echo "🔍 Trying help command..."
            if timeout 10s "$binary" --help 2>&1; then
                echo "✅ $binary - help command works"
            else
                echo "❌ $binary - help command also failed"
                echo "🔍 Checking if it's a dynamic linking issue..."
                ldd "$binary" 2>/dev/null || echo "Binary is statically linked or ldd not available"
            fi
            exit 1
        fi
    fi
    
    echo "---"
done

# Calculate binary sizes
echo "📊 Binary sizes:"
for binary in $(find dist -name "kavach*" -type f); do
    size=$(du -h "$binary" | cut -f1)
    echo "$(basename "$binary"): $size"
done

# Create a summary
echo "📋 Build Summary:"
echo "Version: ${GORELEASER_CURRENT_TAG:-dev}"
echo "Build Time: $(date -u)"
echo "Git Commit: $(git rev-parse HEAD)"
echo "Total Binaries: $(find dist -name "kavach*" -type f | wc -l)"

echo "✅ Post-build process completed successfully!" 