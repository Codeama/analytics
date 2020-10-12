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
export interface AnalyticsProps extends StackProps {
  namespace: string;
}
export class AnalyticsStack extends Stack {
  private role: Role;
  private api: CfnApi;
  private namespace: string;
  private viewsRouteKey: Views;
  private defaultRouteKey: Default;
  private snsTopic: Topic;

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

    this.createQueueHandlers();
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

  createQueueHandlers = () => {
    // POST
    const postHandler = new QueueHandler(this, this.namespace + 'post', {
      name: this.namespace + 'postQueueFunc',
      lambdaDir: './../../analytics-service/post-handler/dist/main.zip',
      topic: this.snsTopic,
    });

    const postSubscriber = postHandler.createSubscriptionFilters(['post_view']);
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
  };
}
