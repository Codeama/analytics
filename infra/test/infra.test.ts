import {
  expect as expectCDK,
  matchTemplate,
  MatchStyle,
} from "@aws-cdk/assert";
import * as cdk from "@aws-cdk/core";
import * as Infra from "../lib/analytics-stack";

test("Empty Stack", () => {
  const app = new cdk.App();
  // WHEN
  const stack = new Infra.AnalyticsStack(app, "MyTestStack");
  // THEN
  expectCDK(stack).to(
    matchTemplate(
      {
        Resources: {},
      },
      MatchStyle.EXACT
    )
  );
});
