# benefits-crawler
n system developed for a code challenge for an interview.

## Running the project
Using docker-compose, the project can be run easily. The default configuration runs Redis, Elasticsearch, RabbitMQ, the API server container, and 3 replicas of the Crawler server:
```
docker-compose up
```

Alternatively, the project can be run by manually running the go files or binaries. The -r option needs to point to the root directory of this repo. The compose-dev.yaml file can be used to run Redis, Elasticsearch and RabbitMQ. As an example:
```
docker-compose -f compose-dev.yaml
go run cmd/api-server/api-server.go -r ./
go run cmd/crawler-server/crawler-server.go -r ./
```
The frontend client can then either be built with `npm run build` (and will be served by the API server) or run in development mode with `npm run start`, in the /client directory.

## Usage
The API server listens to the 8080 port and exposes the `/api/crawlerProcesses` endpoint. This endpoint is used to create and check crawler processes. Using POST, the API creates a crawler process that is then processed by the Crawler servers. The information and status of the process can be retrieved by using GET and the id that was returned in the body from POST. The actual result can be retrieved by calling GET `api/benefits/{cpf}`. Here is an example of curl calls:
```
curl -XPOST -H "Content-type: application/json" -d '{
    "cpf": "000.000.000-00",
    "username": "user",
    "password": "pass"
}' 'localhost:8080/api/crawlerProcesses'
```
Which returns:
```
{
    "_id": "j3eu0IkB7cxQHVshlVKl",
    "cpf": "000.000.000-00",
    "username": "user",
    "password": "pass"
    "process_state": "Created"
}
```

Using GET:
```
curl -XGET 'localhost:8080/api/crawlerProcesses/j3eu0IkB7cxQHVshlVKl'
```
Considering the process was successful it should return:
```
{
    "_id": "j3eu0IkB7cxQHVshlVKl",
    "cpf": "000.000.000-00",
    "username": "user",
    "password": "pass"
    "process_state": "Done"
}
```
The `process_state` can return `Created`, `Running`, `Canceled`, `Failed`, and `Done`. `Canceled` is returned in case there is a result cached from a previous request and `Failed` is returned in case the crawler is not able to retrieve the information.

The results can be retrieved with the API:
```
curl -XGET 'localhost:8080/api/benefits/000.000.000-00'
```
The results can also be checked in the frontend client, which should be accessible by `localhost:8080` or, in case it's running in development mode `localhost:3000`

## Project structure and explanation
This project was built using Go for both the API and Crawler servers. They have shared code in the infra, models, and config directories. Their main files are both in the cmd directory.

The API server uses the GIN framework to serve the static frontend files and the API endpoints. Most of the code is located in the api folder.

The Crawler server uses the colly library for initial data scraping, but most of the crawling is done with the rod library, which runs a Chromium instance to retrieve the information. Most of the code is located in the consumer and crawler directories.

When the API receives a POST request in `/api/crawlerProcesses`, it indexes a document that represents the crawler's process and creates a message in RabbitMQ. This message can be consumed by any crawler server that is consuming the queue. A crawler server checks if a CPF result was already cached in Redis and if not, starts crawling. If successful, the results are indexed in Elasticsearch, or the indexed document is updated, and the result is put in the Redis cache.

The frontend client was created using create-react-app and the tailwindcss, because it was a very simple and small application.
