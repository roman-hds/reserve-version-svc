A service is accessible at port 8080

HTTP Endpoints
/live: 
Description: A liveness endpoint is used to indicate whether the service is to be restarted. (HTTP 200 if healthy, HTTP 503 if unhealthy)
Sample Usage: curl http://localhost:8080/live

/version: 
Description: version endpoint returns an incremented build number and creates a new directory in Artifactory with that number. Accepts 'branch' parameter
Sample Usage: curl http://localhost:8080/version\?branch=reserve
