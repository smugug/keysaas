import fs from 'fs';

export const K8S_API_HOST = '192.168.49.2';
export const K8S_API_PORT = '8443';
export const NAMESPACE = 'customer2';
export const KEYSAAS_API_GROUP = 'keysaascontroller.keysaas';
export const KEY_SAAS_API_VERSION = 'v1';

export const K8S_BASE_IP_URL = `https://${K8S_API_HOST}:${K8S_API_PORT}/apis/${KEYSAAS_API_GROUP}/${KEY_SAAS_API_VERSION}/namespaces/${NAMESPACE}/keysaases`;

export const getK8sToken = (): string => {
    return fs.readFileSync('token', 'utf8');
};