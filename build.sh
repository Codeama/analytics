#!/bin/bash
# TODO Go package mgmt malarkey >.<

# with args ---> for local development
# if [[ -n "$1" ]]; then
#     dirname=$1
#     cd "$dirname" || exit
#     GOOS=linux go build -o dist/main main.go
#     zip main.zip main
#     #  upload/update function code
# fi

service_handlers="
services/default_route
services/connect
services/copy
services/home
services/posts
services/views
"

# Build and package Go Lambda binaries
for dir in $service_handlers
do
cd "$dir" || exit
echo "Building: ${dir} lambda"
go mod tidy
GOOS=linux go build -o dist/main cmd/main.go
zip -r -j dist/main.zip dist/main
cd - || exit 
done

exit 0
