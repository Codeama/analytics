#!/bin/bash
# TODO Package mgmt malarkey >:<
# TODO take lambda-directory-name aas arg and run a different function for local devvelopment
service_handlers="
analytics-service/default
analytics-service/views
analytics-service/post-handler
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
