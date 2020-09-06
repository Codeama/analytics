#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from '@aws-cdk/core';
import { AnalyticsStack } from '../lib/analytics-stack';

const app = new cdk.App();
new AnalyticsStack(app, 'Analytics');
