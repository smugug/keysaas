import axios from 'axios';
import {BASE_URL, PROMETHEUS_URL, THEMES_URL} from './constants';

export const getKeySaaSInstances = async () => {
  try {
    const response = await axios.get(BASE_URL);
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const getKeySaaSInstance = async (name: string) => {
  try {
    const url = `${BASE_URL}/${name}`;
    const response = await axios.get(url);
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const createKeySaaSInstance = async (keySaaS: any) => {
  try {
    const response = await axios.post(BASE_URL, keySaaS);
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const updateKeySaaSInstance = async (name: string, keySaaS: any) => {
  try {
    const url = `${BASE_URL}/${name}`;
    const response = await axios.put(url, keySaaS);
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const deleteKeySaaSInstance = async (name: string) => {
  try {
    const url = `${BASE_URL}/${name}`;
    const response = await axios.delete(url);
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const getPrometheusMetrics = async (query: string) => {
  try {
    const url = `${PROMETHEUS_URL}`;
    const response = await axios.get(url, {
      params: { query },
    });
    return response.data;
  } catch (error) {
    throw error;
  }
}

export const uploadTheme = async (name: string, formData: FormData) => {
  try {
    const response = await axios.post(`${THEMES_URL}/${name}`, FormData);
    return response.data
  } catch (error) {
    throw error;
  }
}