#!/bin/bash

cd infra || exit
npm i

# Deploys the CDK stacks as they're named in infra.ts
cdk deploy "${NAMESPACE}AnalyticsDataStore"
cdk deploy "${NAMESPACE}AnalyticsApi"
now=$(date +"%x - %T")
echo "Last run : $now"