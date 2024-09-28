export const K8S_API_HOST = '192.168.49.2';
export const K8S_API_PORT = '8443';
export const NAMESPACE = 'customer2';
export const KEYSAAS_API_GROUP = 'keysaascontroller.keysaas';
export const KEY_SAAS_API_VERSION = 'v1';
// remove https://${K8S_API_HOST}:${K8S_API_PORT}, add proxy to package.json to bypass cors (development stage)
export const BASE_URL = `/apis/${KEYSAAS_API_GROUP}/${KEY_SAAS_API_VERSION}/namespaces/${NAMESPACE}/keysaases`;

export const getToken = (): string | null => {
    return localStorage.getItem('token');
};