# smart-url 

Smart-url is a simple url shortening service. 

## How to Run

**For running in local development environment**:
  ```bash
  make local-run
  ```
**For running in docker local development environment**:
  ```bash
  make docker-run
  ```

## Sample Curls
 ```bash
  curl --location --request POST 'localhost:8080/shorten' \
--header 'Content-Type: application/json' \
--data '{
    "url":"https://www.google.com"
}'

{"short_url":"/XyClUW"}

curl --location --request GET 'localhost:8080/XyClUW' 


curl --location --request GET 'localhost:8080/metrics' 
  ```