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
analytics-service/default
analytics-service/views
analytics-service/home-handler
analytics-service/post-handler
analytics-service/profile-handler
"

# Build and package Go Lambda binaries
for dir in $service_handlers
do
cd "$dir" || exit
lambda_name="$(cut -d'/' -f2 <<<"$dir")" # Splits and retrieves filename
echo "Building: ${lambda_name} lambda"
GOOS=linux go build -o dist/main main.go
zip -r -j dist/main.zip dist/main
cd - || exit 
done

exit 0
