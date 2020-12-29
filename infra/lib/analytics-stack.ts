import {
  Construct,
  ConcreteDependable,
  StackProps,
  Stack,
} from '@aws-cdk/core';
import { Role, ServicePrincipal } from '@aws-cdk/aws-iam';
import { CfnApi, CfnDeployment, CfnStage } from '@aws-cdk/aws-apigatewayv2';
import { Topic } from '@aws-cdk/aws-sns';
import { Default } from './routes/default';
import { Views } from './routes/views';
import { lambdaPolicy } from './policies';
import { HitsHandler } from './subscriber';
import { config } from './config';

export interface AnalyticsProps extends StackProps {
  namespace: string;
}
export class AnalyticsStack extends Stack {
  private role: Role;
  private api: CfnApi;
  private namespace: string;
  private viewsRouteKey: Views;
  private defaultRouteKey: Default;
  // The SNS topic for hit counter lambdas to subscribe to
  private snsTopic: Topic;
  private homeHitsHandler: HitsHandler;
  private profileHitsHandler: HitsHandler;
  private postHitsHandler: HitsHandler;

  constructor(scope: Construct, id: string, props: AnalyticsProps) {
    super(scope, id, props);

    this.namespace = props.namespace;
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

    // This is the URL that gets generated after successful deployment of the API
    // where stage name (see below) is this.namespace
    const url = `https://${this.api.ref}.execute-api.${config.AWS_REGION}.amazonaws.com/${this.namespace}`;
    this.viewsRouteKey = new Views(this, 'Views', {
      api: this.api,
      role: this.role,
      topic: this.snsTopic,
      topicRegion: config.AWS_REGION as string,
      tableName: config.POST_TABLE_READER,
      apiUrl: url,
      tablePermission: true,
    });

    this.defaultRouteKey = new Default(this, 'Default', {
      api: this.api,
      role: this.role,
    });

    const policy = lambdaPolicy([
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
      autoDeploy: true,
      deploymentId: deployment.ref,
      stageName: this.namespace,
    });

    const dependencies = new ConcreteDependable();
    dependencies.add(this.viewsRouteKey.route);
    dependencies.add(this.defaultRouteKey.route);
    deployment.node.addDependency(dependencies);
  };

  createHitHandlers = () => {
    // HOME
    this.homeHitsHandler = new HitsHandler(this, this.namespace + 'homepage', {
      name: this.namespace + 'homeQueueFunc',
      lambdaDir: './../../analytics-service/home-hits/dist/main.zip',
      topic: this.snsTopic,
      tableName: config.HOME_AND_PROFILE,
      region: config.AWS_REGION as string,
      tablePermission: true,
    });

    this.snsTopic.addSubscription(
      this.homeHitsHandler.createSubscriptionFilters(['homepage_view'])
    );

    // POST
    this.postHitsHandler = new HitsHandler(this, this.namespace + 'post', {
      name: this.namespace + 'postQueueFunc',
      lambdaDir: './../../analytics-service/post-hits/dist/main.zip',
      topic: this.snsTopic,
      tableName: config.POST_TABLE_WRITER,
      region: config.AWS_REGION as string,
      tablePermission: true,
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
        lambdaDir: './../../analytics-service/profile-hits/dist/main.zip',
        topic: this.snsTopic,
        region: config.AWS_REGION as string,
        tableName: config.HOME_AND_PROFILE,
        tablePermission: true,
      }
    );

    const profileSubscriber = this.profileHitsHandler.createSubscriptionFilters(
      ['profile_view']
    );
    this.snsTopic.addSubscription(profileSubscriber);
  };
}
