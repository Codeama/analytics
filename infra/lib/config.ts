const domainConfig = (envName: string) => {
  switch (envName) {
    case 'prod':
      return 'https://bukola.info';
    case 'stage':
      return 'https://staging.bukola.info';
    default:
      return 'http://localhost:8000';
  }
};

export const config = {
  AWS_REGION: process.env.CDK_DEFAULT_REGION,
  POST_TABLE_WRITER: process.env.NAMESPACE + 'PostCountWriter',
  POST_TABLE_READER: process.env.NAMESPACE + 'PostCountReader',
  HOME_AND_PROFILE: process.env.NAMESPACE + 'HomeAndProfile',
  REFERRER_TABLE: process.env.NAMESPACE + 'Referrer',
  DOMAIN_NAME: domainConfig(process.env.NAMESPACE!),
};
