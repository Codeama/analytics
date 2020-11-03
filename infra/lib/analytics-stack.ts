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
import { QueueHandler } from './subscriber';
import { config } from './config';
import { StreamViewType } from '@aws-cdk/aws-dynamodb';
import { Store } from './data-store/table';
import { StreamHandler } from './data-store/db-stream-lambda';
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
  private postCounter: QueueHandler;
  private streamHandler: StreamHandler;

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

    this.viewsRouteKey = new Views(this, 'Views', {
      api: this.api,
      role: this.role,
      topic: this.snsTopic,
      topicRegion: config.AWS_REGION as string,
    });

    this.defaultRouteKey = new Default(this, 'Default', {
      api: this.api,
      role: this.role,
    });

    const policy = lambdaPolicy([
      this.viewsRouteKey.getLambdaArn(),
      this.defaultRouteKey.getLambdaArn(),
    ]);

    this.role.addToPolicy(policy);

    this.createDeployment();

    this.createHitHandlers();

    this.createStorageTables();
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
    dependencies.add(this.viewsRouteKey.getRoute());
    dependencies.add(this.defaultRouteKey.getRoute());
    deployment.node.addDependency(dependencies);
  };

  createHitHandlers = () => {
    // HOME
    const homeHandler = new QueueHandler(this, this.namespace + 'homepage', {
      name: this.namespace + 'homeQueueFunc',
      lambdaDir: './../../analytics-service/home-handler/dist/main.zip',
      topic: this.snsTopic,
    });

    const homeSubscriber = homeHandler.createSubscriptionFilters([
      'homepage_view',
    ]);
    this.snsTopic.addSubscription(homeSubscriber);

    // POST
    this.postCounter = new QueueHandler(this, this.namespace + 'post', {
      name: this.namespace + 'postQueueFunc',
      lambdaDir: './../../analytics-service/post-handler/dist/main.zip',
      topic: this.snsTopic,
    });

    const postSubscriber = this.postCounter.createSubscriptionFilters([
      'post_view',
    ]);
    this.snsTopic.addSubscription(postSubscriber);

    // PROFILE
    const profileHandler = new QueueHandler(this, this.namespace + 'profile', {
      name: this.namespace + 'profileQueueFunc',
      lambdaDir: './../../analytics-service/profile-handler/dist/main.zip',
      topic: this.snsTopic,
    });

    const profileSubscriber = profileHandler.createSubscriptionFilters([
      'profile_view',
    ]);
    this.snsTopic.addSubscription(profileSubscriber);

    // DynamoDB stream lambda function
    this.streamHandler = new StreamHandler(
      this,
      this.namespace + 'StreamLambda',
      {
        lambdaDir: './../../../analytics-service/dynamo-stream/dist/main.zip',
        tableName: 'StreamTable',
      }
    );
  };

  createStorageTables = () => {
    const table = new Store(this, this.namespace + 'Table', {
      tableName: 'PostCountWriter',
      indexName: 'articleId',
      lambdaGrantee: this.postCounter.subscribeFunc,
      stream: StreamViewType.NEW_IMAGE,
    });

    const readerTable = new Store(this, this.namespace + 'ReaderTable', {
      tableName: 'StreamTable',
      indexName: 'articleId',
      lambdaGrantee: this.streamHandler.lambda,
    });
  };
}
