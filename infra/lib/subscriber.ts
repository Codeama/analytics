import * as path from 'path';
import { Function, Runtime, Code } from '@aws-cdk/aws-lambda';
import { Topic, SubscriptionFilter } from '@aws-cdk/aws-sns';
import { Queue } from '@aws-cdk/aws-sqs';
import { SqsSubscription } from '@aws-cdk/aws-sns-subscriptions';
import { SqsEventSource } from '@aws-cdk/aws-lambda-event-sources';
import { Construct } from '@aws-cdk/core';

interface HandlerProps {
  name: string;
  lambdaDir: string;
  topic: Topic;
}

export class QueueHandler extends Construct {
  private subscribeFunc: Function;
  private queue: Queue;
  constructor(scope: Construct, id: string, props: HandlerProps) {
    super(scope, id);
    // TODO increase lambda timeout
    this.subscribeFunc = new Function(this, id + 'Subscriber', {
      runtime: Runtime.GO_1_X,
      code: Code.fromAsset(path.join(__dirname, props.lambdaDir)),
      handler: 'main',
      reservedConcurrentExecutions: 1,
      environment: {
        TOPIC_ARN: props.topic.topicArn,
      },
    });

    this.queue = new Queue(this, id + 'Queue', {
      queueName: id + 'Queue',
    });

    // Lambda SQS trigger
    this.subscribeFunc.addEventSource(
      new SqsEventSource(this.queue, {
        batchSize: 1,
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
