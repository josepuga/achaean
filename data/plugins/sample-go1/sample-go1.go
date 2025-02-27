// Plugin tcp-scan
package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
    // These values are created by Achaean
	pluginID := os.Getenv("PLUGIN_ID")
	progressPipeFile := os.Getenv("PROGRESS_PIPE")

	fmt.Printf("Hello from plugin %s\n", pluginID)
    fmt.Printf("Parameters: %s\n", os.Args[1:])

	// Using a named pipe for the progress bar. https://en.wikipedia.org/wiki/Named_pipe
	progressPipe, err := os.OpenFile(progressPipeFile, os.O_WRONLY, os.ModeNamedPipe)
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
