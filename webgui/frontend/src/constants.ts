export const BASE_URL = `http://kubernetes.local:30005/api/keysaas`;
export const PROMETHEUS_URL = 'http://kubernetes.local:30005/api/prometheus'
export const THEMES_URL = 'http://kubernetes.local:30005/api/themes'
export const getToken = (): string | null => {
    return localStorage.getItem('token');
};
export interface KeySaaSResource {
    apiVersion: string;
    kind: string;
    metadata: {
      name: string;
      namespace: string;
    };
    spec: {
      domainName: string;
      keysaasPassword: string;
      keysaasUsername: string;
      limitsCpu: string;
      limitsMemory: string;
      tls: string;
    };
  }