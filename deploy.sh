#!/bin/bash

cd infra || exit
npm i
cdk bootstrap

# Deploys the CDK stacks as they're named in infra.ts
cdk deploy "${NAMESPACE}AnalyticsDataStore" --require-approval never
cdk deploy "${NAMESPACE}AnalyticsApi" --require-approval never
now=$(date +"%x - %T")
echo "Last run : $now"