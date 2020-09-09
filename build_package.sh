#!/bin/bash

service_handlers="
analytics-service/default
analytics-service/views
"

# Build and package Go Lambda binaries
for dir in $service_handlers
do
cd "$dir" || exit
GOOS=linux go build main.go
zip main.zip main
cd - || exit 
done

exit 0
