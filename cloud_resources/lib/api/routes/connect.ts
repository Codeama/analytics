import * as path from 'path';
import { Function, Runtime, Code } from 'aws-cdk-lib/aws-lambda';

import { config } from '../../../config';
import { CfnRefElement } from 'aws-cdk-lib';
import { Role } from 'aws-cdk-lib/aws-iam';
import { CfnIntegration, CfnRouteResponse, CfnIntegrationResponse, CfnRoute } from 'aws-cdk-lib/aws-apigatewayv2';
import { Construct } from 'constructs';

interface ConnectionProps {
  api: CfnRefElement;
  role: Role;
  domainName: string;
}
export class Connection extends Construct {
  readonly connectFunc: Function;
  readonly route: CfnRoute;
  private lambdaIntegration: CfnIntegration;

  constructor(scope: Construct, id: string, props: ConnectionProps) {
    super(scope, id);

    this.connectFunc = new Function(this, 'ConnectFunction', {
      runtime: Runtime.GO_1_X,
      code: Code.fromAsset(
        path.join(
          __dirname,
          '../../../../analytics-service/connect/dist/main.zip'
        )
      ),
      handler: 'main',
      environment: {
        DOMAIN_NAME: props.domainName,
      },
    });

    this.lambdaIntegration = new CfnIntegration(
      this,
      'ConnectLambdaIntegration',
      {
        apiId: props.api.ref,
        integrationType: 'AWS_PROXY',
        integrationUri: `arn:aws:apigateway:${config.AWS_REGION}:lambda:path/2015-03-31/functions/${this.connectFunc.functionArn}/invocations`,
        credentialsArn: props.role.roleArn, //permission to invoke lambda backend
      }
    );

    this.route = new CfnRoute(this, 'Route', {
      apiId: props.api.ref,
      routeKey: '$connect',
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
}
