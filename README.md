# saiServices

##saiAuth

curl --location --request GET 'http://localhost:8800/register' \
--header 'Token: SomeToken' \
--header 'Content-Type: application/json' \
--data-raw '{"key":"user","password":"12345"}'

curl --location --request GET 'http://localhost:8800/login' \
--header 'Token: SomeToken' \
--header 'Content-Type: application/json' \
--data-raw '{"key":"user","password":"12345"}'
