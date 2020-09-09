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

    /**Custom Lambda for views route key
     * It has a default managed LambdaExecutionRole policy
     */
    const viewsHandler = new Function(this, 'ViewsFunction', {
      runtime: Runtime.GO_1_X,
      code: Code.fromAsset(path.join(__dirname, '../../analytics-service/views/main.zip')),
      handler: 'main',
      // role:
    });

    const defaultHandler = new Function(this, 'DefaultHandler', {
        runtime: Runtime.GO_1_X,
        code: Code.fromAsset(path.join(__dirname, '../../analytics-service/default/main.zip')),
        handler: 'main',
        // role:
      });

    const apiPolicy = new PolicyStatement({
      effect: Effect.ALLOW,
      resources: [
        defaultHandler.functionArn,
        viewsHandler.functionArn,
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

    const viewsIntegration = new CfnIntegration(this, 'Message-Integration', {
      apiId: api.ref,
      integrationType: 'AWS_PROXY',
      integrationUri: 
        `arn:aws:apigateway:${config.AWS_REGION}:lambda:path/2015-03-31/functions/${viewsHandler.functionArn}/invocations`,
      credentialsArn: role.roleArn,
    });

    const defaultIntegration = new CfnIntegration(this, 'Default-Integration', {
        apiId: api.ref,
        integrationType: 'AWS_PROXY',
        integrationUri: 
          `arn:aws:apigateway:${config.AWS_REGION}:lambda:path/2015-03-31/functions/${defaultHandler.functionArn}/invocations`,
        credentialsArn: role.roleArn,
      });

    const viewsRoute = new CfnRoute(this, id + 'viewsRoute', {
      apiId: api.ref,
      routeKey: 'views',
      target: `integrations/${viewsIntegration.ref}`,
      authorizationType: 'NONE',
      routeResponseSelectionExpression: '$default'
    });

    new CfnRouteResponse(this, 'msgResponse', {
        apiId: api.ref,
        routeId: viewsRoute.ref,
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
      integrationId: viewsIntegration.ref,
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
    dependencies.add(viewsRoute);
    dependencies.add(defaultRoute);
    deployment.node.addDependency(dependencies);
  }
}
