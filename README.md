# INFO
There is a /csv/:key endpoint exposed which returns a value associated with a key. 
If we want to add a different store we need to just implement the services.CSVLoader interface.

# Step to run the application
The following steps are to be followed in the root directory of the project to run the service.

go build -o main
docker-compose -f docker-compose.yaml up

P.S: The app runs on port 8000.