import {
  Construct,
  ConcreteDependable,
  StackProps,
  Stack,
} from '@aws-cdk/core';
import { Role, ServicePrincipal } from '@aws-cdk/aws-iam';
import { CfnApi, CfnDeployment, CfnStage } from '@aws-cdk/aws-apigatewayv2';
import { Default } from './routes/default';
import { Views } from './routes/views';
import { lambdaPolicy } from './policy_doc';

export interface AnalyticsProps extends StackProps {
  namespace: string;
}
export class AnalyticsStack extends Stack {
  private role: Role;
  private api: CfnApi;
  private namespace: string;

  constructor(scope: Construct, id: string, props: AnalyticsProps) {
    super(scope, id, props);

    this.namespace = props.namespace;
    this.api = new CfnApi(this, id + 'API', {
      name: this.namespace + 'API',
      protocolType: 'WEBSOCKET',
      routeSelectionExpression: '$request.body.message',
    });

    this.role = new Role(this, id + 'WebsocketRole', {
      assumedBy: new ServicePrincipal('apigateway.amazonaws.com'),
    });

    const viewsRouteKey = new Views(this, id + 'Views', {
      api: this.api,
      role: this.role,
    });

    const defaultRouteKey = new Default(this, id + 'Default', {
      api: this.api,
      role: this.role,
    });

    const policy = lambdaPolicy([
      viewsRouteKey.getLambdaArn(),
      defaultRouteKey.getLambdaArn(),
    ]);

    this.role.addToPolicy(policy);

    // todo CREATE deployment function:::deployApi()
    const deployment = new CfnDeployment(this, id + 'deployment', {
      apiId: this.api.ref,
    });

    new CfnStage(this, id + 'stage', {
      apiId: this.api.ref,
      autoDeploy: true,
      deploymentId: deployment.ref,
      stageName: this.namespace,
    });

    const dependencies = new ConcreteDependable();
    dependencies.add(viewsRouteKey.getRoute());
    dependencies.add(defaultRouteKey.getRoute());
    deployment.node.addDependency(dependencies);
  }
}
