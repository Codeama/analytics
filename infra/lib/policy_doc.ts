import { PolicyStatement, Effect } from '@aws-cdk/aws-iam';

export const lambdaPolicy = (lambdasArn: string[]) => {
  return new PolicyStatement({
    effect: Effect.ALLOW,
    resources: lambdasArn,
    actions: ['lambda:InvokeFunction'],
  });
};
