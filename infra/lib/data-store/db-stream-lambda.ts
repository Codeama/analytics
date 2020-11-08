/**This lambda construct reads
 * from DynamoDB stream and writes to a new table
 */
import * as path from 'path';
import { Construct, Duration } from '@aws-cdk/core';
import { Code, Function, Runtime, StartingPosition } from '@aws-cdk/aws-lambda';
import { DynamoEventSource, SqsDlq } from '@aws-cdk/aws-lambda-event-sources';
import { Table } from '@aws-cdk/aws-dynamodb';
import { Queue } from '@aws-cdk/aws-sqs';

interface StreamProps {
  lambdaDir: string;
  tableName: string;
  region: string;
  triggerSource: Table;
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
        TABLE_REGION: props.region,
      },
    });

    // DLQ
    const dlq = new Queue(this, id + 'StreamsDLQ', {
      queueName: id + 'StreamsDLQ',
    });

    this.lambda.addEventSource(
      new DynamoEventSource(props.triggerSource, {
        startingPosition: StartingPosition.TRIM_HORIZON,
        batchSize: 5,
        bisectBatchOnError: true,
        onFailure: new SqsDlq(dlq),
        retryAttempts: 10,
      })
    );
  }
}
