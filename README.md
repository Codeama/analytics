![Build and Deploy](https://github.com/Codeama/analytics/workflows/Build%20and%20Deploy/badge.svg)

## Contents

[1. What](#What)  
[2. Why](#Why)  
[3. How it Works](#How-it-Works)  
[4. Client to Server](#Client-to-Server)  
[5. WebSocket Client](#WebSocket-Client)  
[6. Tech Stack](#Tech-Stack)  
[7. Cloud Resources](#Cloud Resources)  
[8. Data Storage](#Data-Storage)  
[9. Deploying a Stack](Deploying-a-Stack)

## What

Near real-time event-driven analytics server for my blog site.

## Why

- For the love of data :blush:
- Reader privacy - tracks views not users
- It's fun!

### More About _Why_

This project is also a way for me to practise my Go programming skills by building something I actually use.

### How it Works

Architectural diagram can be found [here](docs/analytics-resources-diagram.png).

The server counts the total hit on a page and determines unique views based on the value of `refreshed` as false sent by the client. If this is not supplied, unique views will be calculated the same way as total views. It is up to the client to decide whether a page is refreshed or not to be counted as non-unique.

The WebSocket server also sends back the total views for each article or blog. This can be added to blog pages. (_I haven't added to mine, it requires CSS_ >.< lol\_).

### Client to Server

The following is the JSON data expected by the server:

      {
        "message": "views", // websocket key route; this is mandatory
        "articleId": string,
        "articleTitle": string,
        "previousPage": string, // optional
        "currentPage": string, // optional
        "refreshed": boolean, // optional
        "referrer": string
      }

### WebSocket Client

You can use any WebSocket client but ultimately, the tracking/tagging is up to you. (Blog post coming _soon_ on how I set this up in my Gatsby blog site).

### Tech Stack

- Go: for AWS Lambda
- NodeJS/TypeScript & AWS CDK: for Infrastructure as Code
- SNS, SQS, DynamoDB, DynamoDB Streams and API Gateway WebSockets

### Cloud Resources

The `cloud_resources` directory contains the CDK app for the resources that constitute the entire infrastructure of the project on AWS.
The resources include: API Gateway (Websocket), Lambda functions, SQS, SNS and DynamoDB.
The API Gateway resources are backed by Lambda functions that process incoming requests via a websocket and publish them to SNS.
SNS is subscribed to by a number of SQS queues that receive notifications based on specific filters: `post_views`, `profile_views` and so on.
Each SQS queue has a corresponding Lambda function which processes messages from the queue and sends them off to DynamoDB.

### Data Storage

The application cloud-hosted data storage consist of four DynamoDB tables:

- HomeAndProfile: stores hit counts of my blog homepage and contacts page
- PostCountWriter: stores hit counts of my blog posts and is listened to by a DynamoDB stream
- PostCountReader: stores data written to DynamoDB stream which is a copy of updates to the PostWriter table. It also serves the data back to the WebSocket client via a Lambda
- Referrer: stores referrer data if available

## Deploying a stack

### Pre-requisite

- Go and NodeJS
- An AWS account
- AWS CLI
- AWS CDK CLI

### Environment variables

- Set your AWS profile using the AWS CLI
- Choose your stack namespace. Allowed namespaces are: `prod`, `stage` or `local`. The namespace you set will map to the client domain that can talk to the analytics server on AWS. The client url is yours to decide. Mine is currently my blog site and so I use `prod` to deploy and set a prod client for `PROD_CLIENT_URL`. `dev` is ideal for deploying and testing locally. A list of enviroment variables to set can be found in `cloud_resources/env.template`

### Build, Package and Deploy

- From the root of the project run the script: `./build_packages.sh`
- Then run `./deploy.sh` to deploy to your personal AWS account

### Unit tests

- The test script `test.sh` can be used to run all the unit tests in the project for the Go Lambda functions

### More Details

My deployed stack serves my blog site which is built with Gatsby and bootstrapped with an open-source template by [Lumen](https://github.com/alxshelepenok/gatsby-starter-lumen) which has pages and posts. I have however added react context components and the WebSocket API client to stream data my analytics server.
