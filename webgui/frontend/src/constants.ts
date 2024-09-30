export const BASE_URL = `http://kubernetes.local:30005/api/keysaas`;
export const PROMETHEUS_URL = 'http://kubernetes.local:30005/api/prometheus'
export const getToken = (): string | null => {
    return localStorage.getItem('token');
};