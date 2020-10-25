import { Construct } from '@aws-cdk/core';
import { Table, AttributeType, StreamViewType } from '@aws-cdk/aws-dynamodb';

interface StoreProps {
  indexName: string;
  tableName: string;
  stream?: StreamViewType;
}
export class Store extends Construct {
  private table: Table;
  constructor(scope: Construct, id: string, props: StoreProps) {
    super(scope, id);

    this.table = new Table(this, 'id' + Table, {
      partitionKey: {
        name: props.indexName,
        type: AttributeType.STRING,
      },
      tableName: props.tableName,
      stream: props.stream,
    });
  }
}
