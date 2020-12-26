import * as path from 'path';
import { Function, Runtime, Code } from '@aws-cdk/aws-lambda';
import { Construct, CfnRefElement, Fn } from '@aws-cdk/core';
import {
  CfnIntegration,
  CfnRoute,
  CfnRouteResponse,
  CfnIntegrationResponse,
} from '@aws-cdk/aws-apigatewayv2';
import { ManagedPolicy, Role } from '@aws-cdk/aws-iam';
import { Topic } from '@aws-cdk/aws-sns';
import { config } from './../config';
import { ReadWriteDynamoDBTable } from '../policies';

interface ViewsProps {
  api: CfnRefElement;
  role: Role;
  topic: Topic;
  topicRegion: string;
  tableName: string;
  apiUrl: string;
  tablePermission?: boolean;
}
export class Views extends Construct {
  readonly viewsFunc: Function;
  private route: CfnRoute;
  private lambdaIntegration: CfnIntegration;

  constructor(scope: Construct, id: string, props: ViewsProps) {
    super(scope, id);

    this.viewsFunc = new Function(this, 'Function', {
      runtime: Runtime.GO_1_X,
      code: Code.fromAsset(
        path.join(
          __dirname,
          '../../../analytics-service/views-stream/dist/main.zip'
        )
      ),
      handler: 'main',
      environment: {
        TOPIC_ARN: props.topic.topicArn,
        TOPIC_REGION: props.topicRegion,
        TABLE_NAME: props.tableName,
        API_URL: props.apiUrl,
      },
    });

    // Grant API invoke permission to lambda
    this.viewsFunc.role?.addManagedPolicy(
      ManagedPolicy.fromAwsManagedPolicyName('AmazonAPIGatewayInvokeFullAccess')
    );

    // DynamoDB permissions
    const tableArn = Fn.importValue(props.tableName + 'Arn');
    const tablePolicy = ReadWriteDynamoDBTable([tableArn]);
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

  getLambdaArn = () => this.viewsFunc.functionArn;
  getRoute = () => this.route;
}
