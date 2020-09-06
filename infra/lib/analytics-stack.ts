import * as path from 'path';
import { Construct, ConcreteDependable, StackProps, Stack } from '@aws-cdk/core';
import {
  Role,
  ServicePrincipal,
  PolicyStatement,
  Effect,
} from '@aws-cdk/aws-iam';
import {
  CfnApi,
  CfnDeployment,
  CfnIntegration,
  CfnIntegrationResponse,
  CfnRoute,
  CfnRouteResponse,
  CfnStage,
} from '@aws-cdk/aws-apigatewayv2';
import { Function, Runtime, Code } from '@aws-cdk/aws-lambda';
import { config } from './config';

export class AnalyticsStack extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);

    /**Custom Lambda for message route key
     * It has a default managed LambdaExecutionRole policy
     */
    const messageFunc = new Function(this, 'MessageHandler', {
      runtime: Runtime.GO_1_X,
      code: Code.fromAsset(path.join(__dirname, '../../analytics-service/message/main.zip')),
      handler: 'main',
      // role:
    });

    const defaultFunc = new Function(this, 'DefaultHandler', {
        runtime: Runtime.GO_1_X,
        code: Code.fromAsset(path.join(__dirname, '../../analytics-service/default/main.zip')),
        handler: 'main',
        // role:
      });

    const apiPolicy = new PolicyStatement({
      effect: Effect.ALLOW,
      resources: [
        defaultFunc.functionArn,
        messageFunc.functionArn,
      ],
      actions: ['lambda:InvokeFunction'],
    });

    const role = new Role(this, id + 'apiWebsocketRole', {
      assumedBy: new ServicePrincipal('apigateway.amazonaws.com'),
    });
    role.addToPolicy(apiPolicy);

    const api = new CfnApi(this, 'AnalyticsApi', {
      name: 'AnalyticsAPI',
      protocolType: 'WEBSOCKET',
      routeSelectionExpression: '$request.body.message',
    });

    const messageIntegration = new CfnIntegration(this, 'Message-Integration', {
      apiId: api.ref,
      integrationType: 'AWS_PROXY',
      integrationUri: 
        `arn:aws:apigateway:${config.AWS_REGION}:lambda:path/2015-03-31/functions/${messageFunc.functionArn}/invocations`,
      credentialsArn: role.roleArn,
    });

    const defaultIntegration = new CfnIntegration(this, 'Default-Integration', {
        apiId: api.ref,
        integrationType: 'AWS_PROXY',
        integrationUri: 
          `arn:aws:apigateway:${config.AWS_REGION}:lambda:path/2015-03-31/functions/${defaultFunc.functionArn}/invocations`,
        credentialsArn: role.roleArn,
      });

    const messageRoute = new CfnRoute(this, id + 'messageRoute', {
      apiId: api.ref,
      routeKey: 'counter',
      target: `integrations/${messageIntegration.ref}`,
      authorizationType: 'NONE',
      routeResponseSelectionExpression: '$default'
    });

    new CfnRouteResponse(this, 'msgResponse', {
        apiId: api.ref,
        routeId: messageRoute.ref,
        routeResponseKey: '$default'
    });

    const defaultRoute = new CfnRoute(this, id + 'defaultRoute', {
        apiId: api.ref,
        routeKey: '$default',
        target: `integrations/${defaultIntegration.ref}`,
        authorizationType: 'NONE',
        routeResponseSelectionExpression: '$default'
      });

    new CfnRouteResponse(this, id + 'dfltResponse', {
        apiId: api.ref,
        routeId: defaultRoute.ref,
        routeResponseKey: '$default'
      });
    
    new CfnIntegrationResponse(this, 'MessageResponse', {
      apiId: api.ref,
      integrationId: messageIntegration.ref,
      integrationResponseKey: '/200/',
    });

    new CfnIntegrationResponse(this, 'DefaultResponse', {
        apiId: api.ref,
        integrationId: defaultIntegration.ref,
        integrationResponseKey: '/200/',
      });

    const deployment = new CfnDeployment(this, id + 'deployment', {
      apiId: api.ref,
    });

    new CfnStage(this, id + 'stage', {
      apiId: api.ref,
      autoDeploy: true,
      deploymentId: deployment.ref,
      stageName: 'dev',
    });

    const dependencies = new ConcreteDependable();
    dependencies.add(messageRoute);
    dependencies.add(defaultRoute);
    deployment.node.addDependency(dependencies);
  }
}
