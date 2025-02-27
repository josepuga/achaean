#!/bin/bash

set -e

PLUGIN_PATH="$(dirname "$(realpath "$0")")"/plugins
PROJECT_NAME=achaean
if [[ "$1" == "" ]]; then
    USAGE=$(cat <<EOF
USE: $0 <plugin_id>.
What does this script do?
    · Creates the directory $PLUGIN_PATH/plugin_id
    · go mod init $PROJECT_NAME/plugins/plugin_id
    · adds dependencies go.mod
    · Creates a basic main.go
    · Creates a basic config.json
    · Creates an empty README
EOF
    )
    echo "$USAGE"
    exit 1
fi

if [[ ! "$1" =~ ^[a-zA-Z0-9_-]+$ ]]; then
    echo "Error: plugin_id can only contain alphanumeric characters, underscores, or hyphens."
    exit 1
fi


dest_dir="$PLUGIN_PATH/$1"
[ -d "$dest_dir" ] && echo "Error: $dest_dir already exists!" && exit 1

mkdir -p "$dest_dir"
cd "$dest_dir"
go mod init "$PROJECT_NAME/plugins/$1"
touch README
cat <<EOF > config.json
{
    "plugin": {
        "id": "$1",
        "name": "Short Name",
        "desc": "Description here...",
        "version": "",
        "author": "",
        "contact": "",
        "type": "",
        "header": ""
    }
}
EOF

#echo "replace achaean/common => ../../common" >> go.mod
#echo "require achaean/common v0.0.0-00010101000000-000000000000" >> go.mod
echo "replace github.com/josepuga/achaean/common => ../../common" >> go.mod
echo "require github.com/josepuga/achaean/common v0.0.0-00010101000000-000000000000" >> go.mod

cat <<EOF > main.go
// Plugin $1
package main

import (
	"achaean/common"
	"fmt"
	"io"
	"os"
	"time"
)

func main() {
	// Read JSON values from stdin.
	jsonBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Println("Hello from plugin $1. JSON content:")
	fmt.Println(string(jsonBytes[:]))

	// Using a named pipe for the progress bar. https://en.wikipedia.org/wiki/Named_pipe
	progressPipe, err := common.OpenProgressPipe("$1") // "$1" ==> pluginID
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer progressPipe.Close()

	// Dummy progress. 0% to 100% (5% increment).
	for p := 0; p <= 100; p += 5 {
		if p == 50 {
			fmt.Println("Half of the progress...")           // Output to Stdout.
			fmt.Fprintln(os.Stderr, "This is a fake error.") // Output to Stderr.
		}
		if p == 100 {
			fmt.Println("Hasta la vista baby.")
		}
		fmt.Fprintf(progressPipe, "%d\n", p) // Send the current % to the named pipe.
		time.Sleep(250 * time.Millisecond)
	}

	// THAT'S IMPORTANT!!!. Plugins must end up before exit writing "DONE" to the progressPipe.
	fmt.Fprintf(progressPipe, "DONE\n")
}

EOF

echo "Go Plugin skeleton sucessfully created at $dest_dir"


