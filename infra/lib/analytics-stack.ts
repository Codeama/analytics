import {
  Construct,
  ConcreteDependable,
  StackProps,
  Stack,
} from "@aws-cdk/core";
import { Role, ServicePrincipal } from "@aws-cdk/aws-iam";
import { CfnApi, CfnDeployment, CfnStage } from "@aws-cdk/aws-apigatewayv2";
import { Default } from "./routes/default";
import { Views } from "./routes/views";
import { lambdaPolicy } from "./policy_doc";

export class AnalyticsStack extends Stack {
  private role: Role;
  private api: CfnApi;

  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);

    this.api = new CfnApi(this, id + "API", {
      name: "API",
      protocolType: "WEBSOCKET",
      routeSelectionExpression: "$request.body.message",
    });

    this.role = new Role(this, id + "WebsocketRole", {
      assumedBy: new ServicePrincipal("apigateway.amazonaws.com"),
    });

    const viewsRouteKey = new Views(this, id, {
      api: this.api,
      role: this.role,
    });

    const defaultRouteKey = new Default(this, id, {
      api: this.api,
      role: this.role,
    });

    const policy = lambdaPolicy([
      viewsRouteKey.getLambdaArn(),
      defaultRouteKey.getLambdaArn(),
    ]);

    this.role.addToPolicy(policy);

    // todo CREATE deployment function:::deployApi()
    const deployment = new CfnDeployment(this, id + "deployment", {
      apiId: this.api.ref,
    });

    new CfnStage(this, id + "stage", {
      apiId: this.api.ref,
      autoDeploy: true,
      deploymentId: deployment.ref,
      stageName: "dev",
    });

    const dependencies = new ConcreteDependable();
    dependencies.add(viewsRouteKey.getRoute());
    dependencies.add(defaultRouteKey.getRoute());
    deployment.node.addDependency(dependencies);
  }
}
