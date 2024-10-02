import { Router } from 'express';
import request from 'request';
import { token,cert,BASE_KEYSAAS_URL, BASE_DEPLOYMENT_URL } from '../server';
import multer from 'multer';

const router = Router();
const storage = multer.diskStorage({
    destination: (req, file, cb) => {
        cb(null, `/themes/${req.params.name}`);
    },
    filename: (req, file, cb) => {
        cb(null, file.originalname); // Keep original file name
    },
});

// Define the multer upload instance
const upload = multer({ storage });

// Function to get the deployment object from the Kubernetes API
// function getDeployment(deploymentName: string): Promise<any> {
//     return new Promise((resolve, reject) => {
//         const url = `${BASE_DEPLOYMENT_URL}/${deploymentName}`;
//         console.log("get deployment 1")
//         const options = {
//             url: url,
//             headers: {
//                 Authorization: `Bearer ${token}`,
//             },
//             agentOptions: {
//                 ca: cert,
//                 rejectUnauthorized: true,
//             },
//             json: true,
//         };
  
//         request.get(options, (error, response, body) => {
//             if (error) {
//                 return reject(`Error fetching deployment: ${error}`);
//             }
//             if (response.statusCode !== 200) {
//                 return reject(`Failed to get deployment. Status code: ${response.statusCode}`);
//             }
//             console.log("get deployment 2")
//             resolve(body);
//         });
//     });
// }

// function getVolumePath(deployment: any): string {
//     console.log("get volume path")
//     const volumes = deployment.spec.template.spec.volumes;
//     for (const volume of volumes) {
//         if (volume.name==="theme-volume" && volume.persistentVolumeClaim) {
//             return path.join('/tmp/hostpath_pv', deployment.metadata.name); // Adjust the base path as needed
//         }
//     }
//     throw new Error('Persistent volume path not found in deployment');
// }

function restartDeployment(deploymentName: string) { 
    const patchBody = [{"op": "replace", "path": "/metadata/annotations/keysaas~1restart", "value": "1"}]
    console.log("restart deployment 1")
    const options = {
        url: `${BASE_KEYSAAS_URL}/${deploymentName}`,
        headers: {
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/json-patch+json',
        },
        agentOptions: {
            ca: cert, 
            rejectUnauthorized: true,
        },
        body: JSON.stringify(patchBody),
    };
  
    request.patch(options, (error, response, body) => {
        console.log("restart deployment 2")
        if (error) {
            throw new Error(error);
        }
    });
}

router.post('/:name', upload.single('file'), (req, res) => {
    try {
        // Check if file is uploaded
        console.log("START UPLOADING 1")
        if (!req.file) {
            res.status(400).json({ message: 'No file uploaded' });
            return
        }
        console.log("START UPLOADING 2")
        res.status(200).json({ message: 'Theme uploaded successfully'});
        restartDeployment(req.params.name)
    } catch (error) {
        console.error('Error uploading theme:', error);
        res.status(500).json({ message: 'Failed to upload theme', error });
    }
});

export default router