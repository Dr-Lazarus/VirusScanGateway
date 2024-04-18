#!/bin/bash

# The directory containing the files to upload
DIRECTORY="./test"

# The endpoint to which the files will be uploaded
ENDPOINT="http://localhost:8080/upload"

# Loop over each file in the directory
for file in "$DIRECTORY"/*
do
  if [ -f "$file" ]; then
    echo "Uploading $file..."
    curl -X POST -F "file=@$file" $ENDPOINT
    echo "*************** COMPLETED ***************" # Print a new line for readability
  fi
done
