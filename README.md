pgservice - Coding Challenge
=================

Running the app locally
---
Clone the repo and navigate to the root folder
```
git clone git@github.com:rumsrami/pgservice.git
cd pgservice
```

### Using docker-compose

1. Install docker-cli and docker-compose
2. From the root folder run the pgservice.
``` 
docker-compose up --build
```
- > The server will run and listen to requests
- > Maps ports 5000:5000
3. Test using cURL or Postman
- POST:
```
curl --location --request POST 'http://0.0.0.0:5000/post-data/title' \
--header 'Content-Type: application/json' \
--data-raw '{
    "Title": "Name"
}'
```
- GET:
``` 
curl --location --request GET 'http://0.0.0.0:5000/get-data/Name'
```
4. Teardown the created containers and network
```
docker-compose down
``` 