#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from '@aws-cdk/core';
import { AnalyticsStack } from '../lib/analytics-stack';

if (!process.env.NAMESPACE) {
  throw Error('NAMESPACE environment must be set.');
}

const namespace = process.env.NAMESPACE as string;

const app = new cdk.App();
new AnalyticsStack(app, namespace + 'Analytics', {
  namespace,
});
