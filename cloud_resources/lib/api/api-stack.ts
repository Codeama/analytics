
import { Default } from './routes/default';
import { Views } from './routes/views';
import { lambdaPolicy } from './policies';
import { HitsHandler } from './subscriber';
import { config } from '../../config';

import { Connection } from './routes/connect';
import { StackProps, Stack } from 'aws-cdk-lib';
import {CfnDeployment, CfnStage } from 'aws-cdk-lib/aws-apigatewayv2';
import { CfnApi } from 'aws-cdk-lib/aws-apigatewayv2';
import { Role, ServicePrincipal, ManagedPolicy } from 'aws-cdk-lib/aws-iam';
import { LogGroup, RetentionDays } from 'aws-cdk-lib/aws-logs';
import { Topic } from 'aws-cdk-lib/aws-sns';
import { Construct } from 'constructs';
import { CfnAccount } from 'aws-cdk-lib/aws-apigateway';

export interface ApiProps extends StackProps {
  namespace: string;
}
export class ApiStack extends Stack {
  private role: Role;
  private api: CfnApi;
  private namespace: string;
  private connectRouteKey: Connection;
  private viewsRouteKey: Views;
  private defaultRouteKey: Default;
  // The SNS topic for hit counter lambdas to subscribe to
  private snsTopic: Topic;
  private homeHitsHandler: HitsHandler;
  private profileHitsHandler: HitsHandler;
  private postHitsHandler: HitsHandler;
  private apiLogGroup: LogGroup;

  constructor(scope: Construct, id: string, props: ApiProps) {
    super(scope, id, props);

    this.namespace = props.namespace;

    // Creates cloudwatch role for API Gateway

    const cloudwatchRole = new Role(this, id + 'CloudWatchRole', {
      assumedBy: new ServicePrincipal('apigateway.amazonaws.com'),
      managedPolicies: [
        ManagedPolicy.fromAwsManagedPolicyName(
          'service-role/AmazonAPIGatewayPushToCloudWatchLogs'
        ),
      ],
    });

    const enableApiLogging = new CfnAccount(this, id + 'CloudWatch', {
      cloudWatchRoleArn: cloudwatchRole.roleArn,
    });

    this.snsTopic = new Topic(this, id + 'Topic', {
      topicName: this.namespace + id + 'Topic',
    });

    this.api = new CfnApi(this, id + 'API', {
      name: this.namespace + 'API',
      protocolType: 'WEBSOCKET',
      routeSelectionExpression: '$request.body.message',
    });

    this.role = new Role(this, id + 'WebsocketRole', {
      assumedBy: new ServicePrincipal('apigateway.amazonaws.com'),
    });

    this.apiLogGroup = new LogGroup(this, id + 'WebSocketLogGoup', {
      retention: RetentionDays.SIX_MONTHS,
    });

    this.connectRouteKey = new Connection(this, 'Connect', {
      api: this.api,
      role: this.role,
      domainName: config.DOMAIN_NAME as string,
    });

    // This is the URL that gets generated after successful deployment of the API
    // where stage name (see below) is this.namespace
    const url = `https://${this.api.ref}.execute-api.${config.AWS_REGION}.amazonaws.com/${this.namespace}`;
    this.viewsRouteKey = new Views(this, 'Views', {
      api: this.api,
      role: this.role,
      topic: this.snsTopic,
      region: config.AWS_REGION as string,
      postTableName: config.POST_TABLE_READER,
      referrerTableName: config.REFERRER_TABLE,
      connectionUrl: url,
      tablePermission: true,
      domainName: config.DOMAIN_NAME,
    });

    this.defaultRouteKey = new Default(this, 'Default', {
      api: this.api,
      role: this.role,
    });

    const policy = lambdaPolicy([
      this.connectRouteKey.connectFunc.functionArn,
      this.viewsRouteKey.viewsFunc.functionArn,
      this.defaultRouteKey.defaultFunc.functionArn,
    ]);

    this.role.addToPolicy(policy);

    this.createDeployment();

    this.createHitHandlers();
  }

  // Bundles up the API resources for deployment
  createDeployment = () => {
    const deployment = new CfnDeployment(this, this.namespace + 'deployment', {
      apiId: this.api.ref,
    });

    new CfnStage(this, this.namespace + 'stage', {
      apiId: this.api.ref,
      // autoDeploy: true, //changes to Stage constructs will not deploy unless autoDeploy is false
      deploymentId: deployment.ref,
      stageName: this.namespace,
      accessLogSettings: {
        destinationArn: this.apiLogGroup.logGroupArn,
        format:
          '{"requestId":"$context.requestId", "caller-domain":"$context.domainName", "user":"$context.identity.user","requestTime":"$context.requestTime", "eventType":"$context.eventType","routeKey":"$context.routeKey", "status":"$context.status","connectionId":"$context.connectionId"}',
      },
    });

    // const dependencies = new ConcreteDependable();
    // dependencies.add(this.viewsRouteKey.route);
    // dependencies.add(this.defaultRouteKey.route);
    // deployment.node.addDependency(dependencies);
  };

  createHitHandlers = () => {
    // HOME
    this.homeHitsHandler = new HitsHandler(this, this.namespace + 'homepage', {
      name: this.namespace + 'homeQueueFunc',
      lambdaDir: '../../../analytics-service/home-hits/dist/main.zip',
      topic: this.snsTopic,
      tableName: config.HOME_AND_PROFILE,
      region: config.AWS_REGION as string,
      tablePermission: true,
      domainName: config.DOMAIN_NAME,
    });

    this.snsTopic.addSubscription(
      this.homeHitsHandler.createSubscriptionFilters(['homepage_view'])
    );

    // POST
    this.postHitsHandler = new HitsHandler(this, this.namespace + 'post', {
      name: this.namespace + 'postQueueFunc',
      lambdaDir: '../../../analytics-service/post-hits/dist/main.zip',
      topic: this.snsTopic,
      tableName: config.POST_TABLE_WRITER,
      region: config.AWS_REGION as string,
      tablePermission: true,
      domainName: config.DOMAIN_NAME as string,
    });

    this.snsTopic.addSubscription(
      this.postHitsHandler.createSubscriptionFilters(['post_view'])
    );

    // PROFILE
    this.profileHitsHandler = new HitsHandler(
      this,
      this.namespace + 'profile',
      {
        name: this.namespace + 'profileQueueFunc',
        lambdaDir: '../../../analytics-service/profile-hits/dist/main.zip',
        topic: this.snsTopic,
        region: config.AWS_REGION as string,
        tableName: config.HOME_AND_PROFILE,
        tablePermission: true,
        domainName: config.DOMAIN_NAME,
      }
    );

    const profileSubscriber = this.profileHitsHandler.createSubscriptionFilters(
      ['profile_view']
    );
    this.snsTopic.addSubscription(profileSubscriber);
  };
}
