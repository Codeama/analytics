![Build and Deploy](https://github.com/Codeama/analytics/workflows/Build%20and%20Deploy/badge.svg)

## What

Event driven near real-time analytics server for my blog site at https://bukola.info/.

## Why

- For the love of data :blush:
- Reader privacy - tracks views not users
- It's fun!

## Technology used

## Project structure

### Infrastructure

The `infra` directory contains the CDK app for the resources that constitute the entire infrastructure of the project on AWS.
The resources include: API Gateway (Websocket), Lambda functions, SQS, SNS and DynamoDB.
The API Gateway resources are backed by Lambda functions that process incoming requests via a websocket and publish them to SNS.
SNS is subscribed to by a number of SQS queues that receive notifications based on specific filters: `post_views`, `profile_views` and so on.
Each SQS queue has a corresponding Lambda function which processes messages from the queue and sends them off to DynamoDB.

## Deploying a stack

### Pre-requisite

- AWS CLI
- AWS CDK CLI
  Set your AWS profile using the AWS CLI
  Set your stack namespace: `export NAMESPACE=my-stack`
  From the root of this project run the script: `./build_packages.sh`
