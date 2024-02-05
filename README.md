# go-aws-serverless-task-management

A serverless REST API to manage projects and tasks using Golang

### Building the zip file manually to upload to AWS

The first way I used to deploy the Lambda code was manually creating the zip file with the binary file inside.

To do that, the first step is to run the command below from the root path to build the binary file called main inside the build folder:

`GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/main -tags lambda.norpc cmd/main.go`

The second step is to create the main.zip file inside the same build folder, with the binary inside it.

And it's done!

## Data Model

### Why did I choose the DynamoDB as database for this project

I learned that DynamoDB is a good choice when working with larges amount of data and when latency is an important requirement, even when the data size keep growing. This app was created on study purposes, so why did I choose this database? As I'm studying and preparing for different usecases in real-life scenarios, I wanted to explore the challenge of modeling the data in the unique Dynamo NoSQL approach.

### Cons of using DynamoDB

If the application keeps changing it's data model, DynamoDB can be a bad choice, because it is better when you already know your data access patterns so you can model the data accordingly. You define the partition and sort key based on these patterns. If your application evolves over time, it's difficult to adapt the data model to the new requirements.

In this case, I've based the data model on data access patterns like "find all projects that a user is part of". To do so, I needed to create a partition key with the username and the sort key with the project id, so I can combine these keys to find these projects. If I need to to the reverse query (find all users that are part of a project), I can create a GSI (Global Secondary Index) to switch the sort key and the partition key, and then I can do the proper query.

## Unit Testing

To simply run the unit tests, run `go test ./...`.

To see the coverage, firstly run `go test -coverprofile=c.out ./...` and then run `go tool cover -html="c.out" ./...`.
