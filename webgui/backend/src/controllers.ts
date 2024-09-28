import { Request, Response } from 'express';
import axios from 'axios';
import './constants'
import { getK8sToken, K8S_BASE_IP_URL } from './constants';



// List Keysaas resources
export const listKeysaasResources = async (req: Request, res: Response) => {
    try {
        const response = await axios.get(K8S_BASE_IP_URL, {
            headers: {
                Authorization: `Bearer ${getK8sToken()}`,
                'Content-Type': 'application/json',
            },
            httpsAgent: new (require('https').Agent)({ rejectUnauthorized: false }) // Disable SSL verification for simplicity; use CA certs in production
        });
        res.json(response.data);
    } catch (error) {
        console.error('Error listing Keysaas resources:', error);
        res.status(500).json({ error: 'Failed to list Keysaas resources' });
    }
};

// Create a Keysaas resource
export const createKeysaasResource = async (req: Request, res: Response) => {
    const keysaasSpec = req.body;

    try {
        const response = await axios.post(K8S_BASE_IP_URL, keysaasSpec, {
            headers: {
                Authorization: `Bearer ${getK8sToken()}`,
                'Content-Type': 'application/json',
            },
            httpsAgent: new (require('https').Agent)({ rejectUnauthorized: false }) // Disable SSL verification for simplicity; use CA certs in production
        });
        res.status(201).json(response.data);
    } catch (error) {
        console.error('Error creating Keysaas resource:', error);
        res.status(500).json({ error: 'Failed to create Keysaas resource' });
    }
};

// Delete a Keysaas resource by name
export const deleteKeysaasResource = async (req: Request, res: Response) => {
    const { name } = req.params;

    try {
        const response = await axios.delete(`${K8S_BASE_IP_URL}/${name}`, {
            headers: {
                Authorization: `Bearer ${getK8sToken()}`,
                'Content-Type': 'application/json',
            },
            httpsAgent: new (require('https').Agent)({ rejectUnauthorized: false }) // Disable SSL verification for simplicity; use CA certs in production
        });
        res.json(response.data);
    } catch (error) {
        console.error('Error deleting Keysaas resource:', error);
        res.status(500).json({ error: 'Failed to delete Keysaas resource' });
    }
};
