/**
 * @param envName - NAMESPACE that has been set
 * @returns WebSocket client URL that has been set
 */
const domainConfig = (envName: string) => {
  switch (envName) {
    case 'prod':
      return process.env.PROD_CLIENT_URL as string;
    case 'stage':
      return process.env.STAGING_CLIENT_URL as string;
    case 'dev':
      return process.env.LOCAL_CLIENT_URL as string;
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

export const validate = () =>{
  if (!process.env.NAMESPACE) {
    throw Error("NAMESPACE environment variable is required.")
  }
}
