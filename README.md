# quotes-api

This is the readme for the quotes API.

## About the API

The API is a connection to the Database setup by https://github.com/skjaldbaka17/setup-quotes-db which contains around 800.000 quotes, over 25.000 authors and over 10.000 quotes sorted into 12 topics.

The primary motivation for this project was to learn to setup and manage SaaS on AWS, using EC2 and RDS and more.

After less then a month the API was ready but then we decided to scrap it and change it into a Serverless API using aws lambda and API Gateway. Which we did, see repo: https://github.com/Skjaldbaka17/quotel-sls-api/tree/main .

## Requirements

* [Golang](https://golang.org)

## Local dev

### Testing

We use the `testing` package that comes built-in in Golang. Before running the tests you need to create a `.env` file in root with `DATABASE_URL=YOUR_DB_URL` then you can simply run the following commands to test the various functionalities of the api:

For all tests:
```shell
make test
```

For all tests with verbose on:
```shell
make test-verbose
```

For a specific test function
```shell
make test-specific TEST_FUNCTION=<The_function_you_want>
```


### API Documentation

For documenting the API we use Swagger (or OpenAPI) and document each endpoint inside the code with specific comments forexed with `swagger:route`. To compile these comments into a swagger.yaml file you simply run:

```shell
make docs
```

This command will first check if you have the goswagger bin compiled. If it is not installed on your machine the command should install it with the command
```shell
go get -u github.com/go-swagger/go-swagger/cmd/swagger
```

The docs are then hosted at `/docs`.

## Setup process

### Setup EC2

For setting up the EC2 you can use the setupEc2.sh as a guide... Of course it would have been best to do this using cloudformation and just creating an image of the EC2 environment we need and using that image in the cloudformation when spinning up a new instance, but hindsight is 2020. What.

