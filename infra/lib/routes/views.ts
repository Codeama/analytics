import * as path from "path";
import { Function, Runtime, Code } from "@aws-cdk/aws-lambda";
import { Construct, CfnRefElement } from "@aws-cdk/core";
import {
  CfnIntegration,
  CfnRoute,
  CfnRouteResponse,
  CfnIntegrationResponse,
} from "@aws-cdk/aws-apigatewayv2";
import { config } from "./../config";
import { Role } from "@aws-cdk/aws-iam";

interface ViewsProps {
  api: CfnRefElement;
  role: Role;
}
export class Views extends Construct {
  private viewsFunc: Function;
  private route: CfnRoute;
  private lambdaIntegration: CfnIntegration;

  constructor(scope: Construct, id: string, props: ViewsProps) {
    super(scope, id);

    this.viewsFunc = new Function(this, "ViewsFunction", {
      runtime: Runtime.GO_1_X,
      code: Code.fromAsset(
        path.join(__dirname, "../../../analytics-service/views/main.zip")
      ),
      handler: "main",
    });

    this.lambdaIntegration = new CfnIntegration(
      this,
      "ViewsLambdaIntegration",
      {
        apiId: props.api.ref,
        integrationType: "AWS_PROXY",
        integrationUri: `arn:aws:apigateway:${config.AWS_REGION}:lambda:path/2015-03-31/functions/${this.viewsFunc.functionArn}/invocations`,
        credentialsArn: props.role.roleArn,
      }
    );

    this.route = new CfnRoute(this, id + "ViewsRoute", {
      apiId: props.api.ref,
      routeKey: "views",
      target: `integrations/${this.lambdaIntegration.ref}`,
      authorizationType: "NONE",
      routeResponseSelectionExpression: "$default",
    });

    new CfnRouteResponse(this, id + "ViewsRouteResponse", {
      apiId: props.api.ref,
      routeId: this.route.ref,
      routeResponseKey: "$default",
    });

    new CfnIntegrationResponse(this, "ViewsIntegrationResponse", {
      apiId: props.api.ref,
      integrationId: this.lambdaIntegration.ref,
      integrationResponseKey: "/200/",
    });
  }

  getLambdaArn = () => this.viewsFunc.functionArn;
  getRoute = () => this.route;
}