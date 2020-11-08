import * as path from 'path';
import { Function, Runtime, Code } from '@aws-cdk/aws-lambda';
import { Construct, CfnRefElement } from '@aws-cdk/core';
import {
  CfnIntegration,
  CfnRoute,
  CfnRouteResponse,
  CfnIntegrationResponse,
} from '@aws-cdk/aws-apigatewayv2';
import { config } from './../config';
import { Role } from '@aws-cdk/aws-iam';

interface DefaultProps {
  api: CfnRefElement;
  role: Role;
}
export class Default extends Construct {
  private defaultFunc: Function;
  private route: CfnRoute;
  private lambdaIntegration: CfnIntegration;

  constructor(scope: Construct, id: string, props: DefaultProps) {
    super(scope, id);

    this.defaultFunc = new Function(this, 'DefaultFunction', {
      runtime: Runtime.GO_1_X,
      code: Code.fromAsset(
        path.join(
          __dirname,
          '../../../analytics-service/default-stream/dist/main.zip'
        )
      ),
      handler: 'main',
    });

    this.lambdaIntegration = new CfnIntegration(
      this,
      'DefaultLambdaIntegration',
      {
        apiId: props.api.ref,
        integrationType: 'AWS_PROXY',
        integrationUri: `arn:aws:apigateway:${config.AWS_REGION}:lambda:path/2015-03-31/functions/${this.defaultFunc.functionArn}/invocations`,
        credentialsArn: props.role.roleArn,
      }
    );

    this.route = new CfnRoute(this, 'Route', {
      apiId: props.api.ref,
      routeKey: '$default',
      target: `integrations/${this.lambdaIntegration.ref}`,
      authorizationType: 'NONE',
      routeResponseSelectionExpression: '$default',
    });

    new CfnRouteResponse(this, 'RouteResponse', {
      apiId: props.api.ref,
      routeId: this.route.ref,
      routeResponseKey: '$default',
    });

    new CfnIntegrationResponse(this, 'IntegrationResponse', {
      apiId: props.api.ref,
      integrationId: this.lambdaIntegration.ref,
      integrationResponseKey: '/200/',
    });
  }

  getLambdaArn = () => this.defaultFunc.functionArn;
  getRoute = () => this.route;
}
