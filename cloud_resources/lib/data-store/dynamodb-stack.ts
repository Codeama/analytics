
import { StackProps, Stack } from 'aws-cdk-lib';
import { StreamViewType } from 'aws-cdk-lib/aws-dynamodb';
import { Construct } from 'constructs';
import { config } from '../../config';
import { StreamHandler } from './db-stream-lambda';
import { Store } from './table';

export interface DatabaseProps extends StackProps {
  namespace: string;
}

export class DatabaseStack extends Stack {
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

    const streamHandler = new StreamHandler(
      this,
      this.namespace + 'StreamLambda',
      {
        lambdaDir: './../../../services/copy/dist/main.zip',
        tableName: config.POST_TABLE_READER,
        region: config.AWS_REGION as string,
        triggerSource: postWriterTable.table,
      }
    );


    const postReaderTable = new Store(this, this.namespace + 'ReaderTable', {
      tableName: config.POST_TABLE_READER,
      indexName: 'articleId',
      lambdaGrantee: streamHandler.lambda, // grant streamHandler permission to write to this table
    });

    const homeAndProfileTable = new Store(
      this,
      this.namespace + 'HomeAndProfileTable',
      {
        tableName: config.HOME_AND_PROFILE,
        indexName: 'pageName',
      }
    );

    const referrerTable = new Store(this, this.namespace + 'ReferrerTable', {
      tableName: config.REFERRER_TABLE,
      indexName: 'id',
    });
  };
}