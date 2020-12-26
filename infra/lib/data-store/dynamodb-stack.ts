import { StreamViewType } from '@aws-cdk/aws-dynamodb';
import { Stack, Construct, StackProps } from '@aws-cdk/core';
import { config } from '../config';
import { StreamHandler } from './db-stream-lambda';
import { Store } from './table';

export interface DatabaseProps extends StackProps {
  namespace: string;
}

export class DatabaseStack extends Stack {
  private streamHandler: StreamHandler;
  private namespace: string;

  constructor(scope: Construct, id: string, props: DatabaseProps) {
    super(scope, id, props);

    this.namespace = props.namespace;

    this.createStorageTables();
  }

  createStorageTables = () => {
    const postWriterTable = new Store(this, this.namespace + 'WriterTable', {
      tableName: config.POST_TABLE_WRITER,
      indexName: 'articleId',
      stream: StreamViewType.NEW_IMAGE,
    });

    this.streamHandler = new StreamHandler(
      this,
      this.namespace + 'StreamLambda',
      {
        lambdaDir: './../../../analytics-service/dynamo-stream/dist/main.zip',
        tableName: config.POST_TABLE_READER,
        region: config.AWS_REGION as string,
        triggerSource: postWriterTable.table,
        tablePermission: true,
      }
    );

    const postReaderTable = new Store(this, this.namespace + 'ReaderTable', {
      tableName: config.POST_TABLE_READER,
      indexName: 'articleId',
    });
  };
}
