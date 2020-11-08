service_handlers="
analytics-service/default-stream
analytics-service/views-stream
analytics-service/post-hits
analytics-service/profile-hits
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