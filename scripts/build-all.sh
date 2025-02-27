#!/bin/bash
set -e

LDFLAGS="-w -s"
echo "Building Achaea..."
go build -ldflags="$LDFLAGS"

echo "Building Plugins..."
cd plugins/
for d in *; do
    [ ! -d "$d" ] && continue
    # Only Go Plugins...
    if [ ! -f "$d/go.mod" ]; then
        echo "$d - Skipping. It doesn't seem like a Golang plugin."
        continue
    fi
    echo "$d"
    cd "$d"
    go build -ldflags="$LDFLAGS"
    cd ..
done
