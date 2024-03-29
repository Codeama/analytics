import * as path from 'path';
import { Function, Runtime, Code } from '@aws-cdk/aws-lambda';
import { Topic, SubscriptionFilter } from '@aws-cdk/aws-sns';
import { Queue } from '@aws-cdk/aws-sqs';
import { SqsSubscription } from '@aws-cdk/aws-sns-subscriptions';
import { SqsEventSource } from '@aws-cdk/aws-lambda-event-sources';
import { Construct, Duration, Fn } from '@aws-cdk/core';
import { ReadWriteDynamoDBTable } from './policies';

interface HandlerProps {
  name: string;
  lambdaDir: string;
  topic: Topic;
  region: string;
  tableName: string;
  tablePermission: boolean;
  domainName: string;
}

export class HitsHandler extends Construct {
  readonly subscribeFunc: Function;
  private queue: Queue;
  constructor(scope: Construct, id: string, props: HandlerProps) {
    super(scope, id);
    this.subscribeFunc = new Function(this, 'Subscriber', {
      runtime: Runtime.GO_1_X,
      code: Code.fromAsset(path.join(__dirname, props.lambdaDir)),
      handler: 'main',
      reservedConcurrentExecutions: 1, // <--- this is a temp solution to avoid creating a lock mechanism for dynamodb :(
      timeout: Duration.seconds(5),
      environment: {
        TOPIC_ARN: props.topic.topicArn,
        TABLE_REGION: props.region,
        TABLE_NAME: props.tableName,
        DOMAIN_NAME: props.domainName as string,
      },
    });

    // DynamoDB permissions
    const tableArn = Fn.importValue(props.tableName + 'Arn');
    const tablePolicy = ReadWriteDynamoDBTable([tableArn]);
    props.tablePermission
      ? this.subscribeFunc.addToRolePolicy(tablePolicy)
      : null;

    // DLQ
    const dlq = new Queue(this, id + 'DLQ', {
      queueName: id + 'DLQ',
    });

    this.queue = new Queue(this, id + 'Queue', {
      queueName: id + 'Queue',
      deadLetterQueue: {
        queue: dlq,
        maxReceiveCount: 5, // total retries before message lands in dlq
      },
    });

    // Lambda SQS trigger
    this.subscribeFunc.addEventSource(
      new SqsEventSource(this.queue, {
        // batchSize: 3,
      })
    );
  }

  //  no filters
  createSubscription = () => {
    new SqsSubscription(this.queue, {
      rawMessageDelivery: true,
    });
  };

  createSubscriptionFilters = (eventTypes: string[]) => {
    return new SqsSubscription(this.queue, {
      filterPolicy: {
        event_type: SubscriptionFilter.stringFilter({
          whitelist: eventTypes,
        }),
      },
      rawMessageDelivery: true,
    });
  };
}
