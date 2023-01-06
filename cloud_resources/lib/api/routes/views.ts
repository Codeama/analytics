import { CfnRefElement, Fn } from 'aws-cdk-lib';
import { CfnIntegration, CfnRouteResponse, CfnIntegrationResponse, CfnRoute } from 'aws-cdk-lib/aws-apigatewayv2';
import { Role, ManagedPolicy } from 'aws-cdk-lib/aws-iam';
import { Runtime, Code, Function } from 'aws-cdk-lib/aws-lambda';
import { Topic } from 'aws-cdk-lib/aws-sns';
import { Construct } from 'constructs';
import * as path from 'path';

import { config } from '../../../config';
import { ReadWriteDynamoDBTable } from '../policies';

interface ViewsProps {
  api: CfnRefElement;
  role: Role;
  topic: Topic;
  region: string;
  postTableName: string;
  referrerTableName: string;
  connectionUrl: string;
  tablePermission?: boolean;
  domainName: string;
}
export class Views extends Construct {
  readonly viewsFunc: Function;
  readonly route: CfnRoute;
  private lambdaIntegration: CfnIntegration;

  constructor(scope: Construct, id: string, props: ViewsProps) {
    super(scope, id);

    this.viewsFunc = new Function(this, 'Function', {
      runtime: Runtime.GO_1_X,
      code: Code.fromAsset(
        path.join(
          __dirname,
          '../../../../services/views/dist/main.zip'
        )
      ),
      handler: 'main',
      environment: {
        TOPIC_ARN: props.topic.topicArn,
        REGION: props.region,
        POST_TABLE_NAME: props.postTableName,
        REFERRER_TABLE_NAME: props.referrerTableName,
        CONNECTION_URL: props.connectionUrl,
        DOMAIN_NAME: props.domainName,
      },
    });

    // Grant API invoke permission to lambda for calls to the connectionUrl to
    // be able to communicate back to connected clients
    this.viewsFunc.role?.addManagedPolicy(
      ManagedPolicy.fromAwsManagedPolicyName('AmazonAPIGatewayInvokeFullAccess')
    );

    // DynamoDB permissions
    const postTableArn = Fn.importValue(props.postTableName + 'Arn');
    const referrerTableArn = Fn.importValue(props.referrerTableName + 'Arn');
    const tablePolicy = ReadWriteDynamoDBTable([
      postTableArn,
      referrerTableArn,
    ]);
    props.tablePermission ? this.viewsFunc.addToRolePolicy(tablePolicy) : null;

    // Topic permission
    props.topic.grantPublish(this.viewsFunc);

    this.lambdaIntegration = new CfnIntegration(
      this,
      'ViewsLambdaIntegration',
      {
        apiId: props.api.ref,
        integrationType: 'AWS_PROXY',
        integrationUri: `arn:aws:apigateway:${config.AWS_REGION}:lambda:path/2015-03-31/functions/${this.viewsFunc.functionArn}/invocations`,
        credentialsArn: props.role.roleArn,
      }
    );

    this.route = new CfnRoute(this, 'Route', {
      apiId: props.api.ref,
      routeKey: 'views',
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
