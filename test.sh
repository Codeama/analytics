service_handlers="
analytics-service/default
analytics-service/views
analytics-service/post-handler
analytics-service/profile-handler
"

# Run all lambda tests
for dir in $service_handlers
do
cd "$dir" || exit
go test -v -cover ./...
cd - || exit 
done

exit 0