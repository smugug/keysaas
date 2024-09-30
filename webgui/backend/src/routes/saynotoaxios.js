const axios = require('axios');
const fs = require('fs');
const https = require('https');

const token_path =  '/var/run/secrets/kubernetes.io/serviceaccount/token';
const cert_path = '/var/run/secrets/kubernetes.io/serviceaccount/ca.crt'
const readFile = (path) => {
  try {
    const content = fs.readFileSync(path,'utf8');
    return content;
  } catch (error) {
    console.error('Error reading service account', error);
    return '';
  }
};
const token = readFile(token_path);
const cert = readFile(cert_path);
const BASE_URL = 'https://kubernetes.default.svc/apis/keysaascontroller.keysaas/v1/namespaces/customer2/keysaases';

// async function fetchData() {
//   const response = await axios.get(BASE_URL, {
//       headers: { Authorization: `Bearer ${token}}` },
//       agentOptions: new https.Agent({
//         ca: cert,
//       }),
//     });
//   res.json(response);
// }
// fetchData();


const request = require('request');

const options = {
  url: 'https://kubernetes.default.svc/apis/keysaascontroller.keysaas/v1/namespaces/customer2/keysaases',
  headers: {
    'Authorization': `Bearer ${token}`,  // Pass your token here
  },
  agentOptions: {
    ca: cert,  // Use the CA certificate
    rejectUnauthorized: true,  // Verify certificate
  }
};

// Make the GET request
request.get(options, (error, response, body) => {
  if (error) {
    console.error('Request error:', error);
  } else {
    console.log('Response:', body);
  }
});
