/**This lambda construct reads
 * from DynamoDB stream and writes to a new table
 */
import * as path from 'path';
import { Construct, Duration } from '@aws-cdk/core';
import { Code, Function, Runtime } from '@aws-cdk/aws-lambda';

interface StreamProps {
  lambdaDir: string;
  tableName: string;
}
export class StreamHandler extends Construct {
  readonly lambda: Function;
  constructor(scope: Construct, id: string, props: StreamProps) {
    super(scope, id);

    this.lambda = new Function(this, 'StreamHandler', {
      runtime: Runtime.GO_1_X,
      code: Code.fromAsset(path.join(__dirname, props.lambdaDir)),
      handler: 'main',
      reservedConcurrentExecutions: 1, // <--- this is a temp solution to avoid creating a lock mechanism for dynamodb :(
      timeout: Duration.seconds(5),
      environment: {
        TABLE_NAME: props.tableName,
      },
    });
  }
}
