/**Disclosure: This is my attempt at using
 * the builder pattern. At the time of using this pattern,
 * only one table needs DynamoDB Streams enabled.
 * But it makes it easier to modify other tables if need be
 * ....or YAGNI lol
 */

import { Construct } from '@aws-cdk/core';
import { StreamViewType } from '@aws-cdk/aws-dynamodb';
import { Store } from './table';

export interface TableBuilderProps {
  scope: Construct;
  stream: StreamViewType;
  tableName: string;
  indexName: string;
}

export class TableBuilder {
  private name: string;
  private index: string;
  // option to enable stream on the table being created
  private stream?: StreamViewType;
  private props: TableBuilderProps;
  constructor(props: TableBuilderProps) {
    this.props = props;
  }

  setTableName(name: string): TableBuilder {
    this.name = name;
    return this;
  }

  setIndex(index: string): TableBuilder {
    this.index = index;
    return this;
  }

  setStreamType(stream: StreamViewType): TableBuilder {
    this.stream = stream;
    return this;
  }

  createTable(): Store {
    const table = new Store(this.props.scope, this.name, {
      tableName: this.name,
      indexName: this.index,
      stream: this.stream,
    });
    return table;
  }
}
