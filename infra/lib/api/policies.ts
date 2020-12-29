import { PolicyStatement, Effect } from '@aws-cdk/aws-iam';

export const lambdaPolicy = (lambdasArn: string[]) => {
  return new PolicyStatement({
    effect: Effect.ALLOW,
    resources: lambdasArn,
    actions: ['lambda:InvokeFunction'],
  });
};

export const ReadDynamoDBTable = (dynamoDBArn: string[]) => {
  return new PolicyStatement({
    effect: Effect.ALLOW,
    resources: dynamoDBArn,
    actions: ['dynamodb:Query', 'dynamodb:GetItem'],
  });
};

export const ReadWriteDynamoDBTable = (dynamoDBArn: string[]) => {
  return new PolicyStatement({
    effect: Effect.ALLOW,
    resources: dynamoDBArn,
    actions: [
      'dynamodb:PutItem',
      'dynamodb:UpdateItem',
      'dynamodb:DescribeTable',
      'dynamodb:Query',
      'dynamodb:GetItem',
    ],
  });
};
