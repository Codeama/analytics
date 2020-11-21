#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from '@aws-cdk/core';
import { AnalyticsStack } from '../lib/analytics-stack';

if (!process.env.NAMESPACE) {
  throw Error('NAMESPACE environment must be set.');
}

const protect = process.env.NAMESPACE === 'prod' ? true : false;

const namespace = process.env.NAMESPACE as string;

const app = new cdk.App();
new AnalyticsStack(app, namespace + 'Analytics', {
  namespace,
  env: {
    account: process.env.CDK_DEFAULT_ACCOUNT,
    region: process.env.CDK_DEFAULT_REGION,
  },
  terminationProtection: protect,
  description: 'Analytics server for my JAMStack blogsite',
});
