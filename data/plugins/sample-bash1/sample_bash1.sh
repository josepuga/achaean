#!/bin/bash

#TODO: Params https://stackoverflow.com/questions/192249/how-do-i-parse-command-line-arguments-in-bash

# Set bash script location as current directory. Not useful in this sample but...
cd "$(dirname "$(realpath "$0")")" || exit 1

# PROGRESS_PIPE Pipe named for progress. Has been created by Achaean.
# >&2 redirects to stderr
[ ! -p "$PROGRESS_PIPE" ] && >&2 echo "Named pipe $PROGRESS_PIPE not created"  && exit 1

# PLUGIN_ID is an enviroment variable created with the ID set in plugin.json
echo "Hello from plugin $PLUGIN_ID."
echo "Parameters: $*"

# Testing Progress Bar. Sending output to the Progress pipe named.
for ((i=0;i<=100;i=i+5)); 
do 
   echo $i > "$PROGRESS_PIPE"
   sleep 0.25
   if [ $i -eq 50 ]; then 
    >&2  echo "This is a fake error"
    echo "Half of the progress..."
   fi
done
echo "Hasta la vista baby."
sleep 0.50
# DONE, tells Achaean that the script has been finished.
echo "DONE" > "$PROGRESS_PIPE"
 