import axios from 'axios';
import {BASE_URL, getToken} from './constants';

export const getKeySaaSInstances = async () => {
  try {
    const response = await axios.get(BASE_URL, {
      headers: { 'Authorization': `Bearer ${getToken()}` },
    });
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const createKeySaaSInstance = async (keySaaS: any) => {
  try {
    const response = await axios.post(BASE_URL, keySaaS, {
      headers: { 'Authorization': `Bearer ${getToken()}` },
    });
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const updateKeySaaSInstance = async (name: string, keySaaS: any) => {
  try {
    const url = `${BASE_URL}/${name}`;
    const response = await axios.put(url, keySaaS, {
      headers: { 'Authorization': `Bearer ${getToken()}` },
    });
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const deleteKeySaaSInstance = async (name: string) => {
  try {
    const url = `${BASE_URL}/${name}`;
    const response = await axios.delete(url, {
      headers: { 'Authorization': `Bearer ${getToken()}` },
    });
    return response.data;
  } catch (error) {
    throw error;
  }
};
