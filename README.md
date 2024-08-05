# Online Store Backend

Online store backend on golang consists rest-api and gRPC as data exchange. Project based on clean architecture(3 layer of abstraction). 

## API Documentation 
OpenAPI 3 documentation:
`swagger/openapi.json` 

Also you can import Postman collection file from root directory:
`eCom Golang.postman_collection.json`

## Database
As database service used Postgresql. 
### Database scheme
![Alt text](/utils/database/market.jpg? "store database scheme")

# TO-DO
* [x] redis
* [x] product add
* [x] product update
* [x] product delete
* [x] cart module
* [x] categories module
* [ ] statistics module
* [ ] filter for each module
* [x] order module
* [x] brand module
* [x] region module
* [x] user module