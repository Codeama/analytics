import { Construct } from '@aws-cdk/core';
import { Table, AttributeType, StreamViewType } from '@aws-cdk/aws-dynamodb';
import { Function } from '@aws-cdk/aws-lambda';

interface StoreProps {
  indexName: string;
  tableName: string;
  stream?: StreamViewType;
  lambdaGrantee: Function;
}
export class Store extends Construct {
  readonly table: Table;
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

    this.table.grantReadWriteData(props.lambdaGrantee.grantPrincipal);
  }
}
