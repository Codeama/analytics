export const config = {
  AWS_REGION: process.env.CDK_DEFAULT_REGION,
  POST_TABLE_WRITER: 'PostCountWriter',
  POST_TABLE_READER: 'PostCountReader',
};
