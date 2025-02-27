#!/usr/bin/env python3
import sys
import time
import os


def main():
    # These values are created by Achaean
    plugin_id = os.environ["PLUGIN_ID"]
    progress_pipe = os.environ["PROGRESS_PIPE"]
    
    print(f"Hello from plugin {plugin_id}.")
    params_string = " ".join(sys.argv[1:])
    print(f"Parameters: {params_string}")

    # Abrir el named pipe para progreso
    try:
        with open(progress_pipe, 'w', encoding="utf-8") as progress_pipe:
            # Progress Simulation
            for p in range(0, 101, 5):
                if p == 50:
                    print("Half of the progress...")  # Output to stdout
                    print("This is a fake error.", file=sys.stderr)  # Output to stderr

                # Write progress to named pipe
                progress_pipe.write(f"{p}")
                progress_pipe.flush()  # Must be flushed, even with open(.. buffering=1)!!
                if p == 100:
                    print("Hasta la vista baby.")

                time.sleep(0.25)  # A simple delay...

            # Write "DONE" in the progress pipe to tell Achaena "I'm finished"
            progress_pipe.write("DONE ")
            progress_pipe.flush()

    except Exception as e:
        print(f"Error handling progress pipe: {e}", file=sys.stderr)
        sys.stderr.flush()

if __name__ == "__main__":
    main()
