service_handlers="
services/default_route
services/views
services/posts
services/profile
"

# Run all lambda tests
for dir in $service_handlers
do
cd "$dir" || exit
go test -v -cover ./...
# go tool cover -html=c.out -o coverage.html
cd - || exit 
done

exit 0