export const config = {
  AWS_REGION: process.env.CDK_DEFAULT_REGION,
  POST_TABLE_WRITER: process.env.NAMESPACE + 'PostCountWriter',
  POST_TABLE_READER: process.env.NAMESPACE + 'PostCountReader',
};
