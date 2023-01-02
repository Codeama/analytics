/**This lambda construct reads
 * from DynamoDB stream and writes to a new table
 */
import { Duration } from 'aws-cdk-lib';
import { Table } from 'aws-cdk-lib/aws-dynamodb';
import { Runtime, Code, StartingPosition, Function } from 'aws-cdk-lib/aws-lambda';
import { DynamoEventSource, SqsDlq } from 'aws-cdk-lib/aws-lambda-event-sources';
import { Queue } from 'aws-cdk-lib/aws-sqs';
import { Construct } from 'constructs';
import * as path from 'path';


interface StreamProps {
  lambdaDir: string;
  tableName: string;
  region: string;
  triggerSource: Table;
  tablePermission?: boolean;
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
