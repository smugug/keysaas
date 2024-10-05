import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import './KeysaasList.css';
import './ThemeUpload';
import {deleteKeySaaSInstance,getKeySaaSInstances} from '../api'
import ThemeUpload from './ThemeUpload';

export interface KeysaasInstance {
  apiVersion: string; // e.g., "keysaascontroller.keysaas/v1"
  kind: string; // e.g., "Keysaas"
  metadata: {
    annotations: {
      [key: string]: string; // Map of annotations, key-value pairs
    };
    creationTimestamp: string; // e.g., "2024-10-04T09:35:58Z"
    generation: number; // e.g., 2
    name: string; // e.g., "keysaastest"
    namespace: string; // e.g., "customer2"
    resourceVersion: string; // e.g., "185309"
    uid: string; // e.g., "a0edbf6e-1c19-4e06-9f6f-b82f4742b157"
  };
  spec: {
    domainName: string,
    keysaasPassword: string,
    keysaasUsername: string,
    requestsCpu: string,
    requestsMemory: string,
    limitsCpu: string,
    limitsMemory: string,
    scalingThreshold: string,
    minInstances: string,
    maxInstances: string,
    postgresUri: string,
    tls: string,
  };
  status: {
    podName: string; // e.g., "keysaastest-69747c67f5-rglnd"
    secretName: string; // e.g., "keysaastest"
    status: string; // e.g., "Ready"
    url: string; // e.g., "http://keysaastest.kubernetes.local"
  };
}


const KeysaasList: React.FC = () => {
  const [keySaaSInstances, setKeySaaSInstances] = useState<KeysaasInstance[]>([]);
  // const [expanded, setExpanded] = useState<{ [key: string]: boolean }>({});
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  // Function to handle the button click to navigate to the detail page
  const goToDetail = (name: string) => {
    navigate(`/keysaas/${name}`);
  };
  // Fetch Keysaas instances from Kubernetes API
  const fetchKeysaasInstances = async () => {
    setLoading(true);

    

    try {
      const data =await getKeySaaSInstances()
      console.log(data)
      setKeySaaSInstances(data.items);
    } catch (error) {
      setError('Failed to fetch KeySaaS instances');
      console.error(error);
    }
    setLoading(false);
  };

  // Delete Keysaas instance
  const handleDelete = async (name: string) => {
    if (window.confirm(`Are you sure you want to delete ${name}?`)) {
      try {
        await deleteKeySaaSInstance(name)
        // Remove the deleted instance from the list
        setKeySaaSInstances((prevInstances) =>
          prevInstances.filter((instance) => instance.metadata.name !== name)
        );  
      } catch (error) {
        setError(`Failed to delete KeySaaS instance: ${name}`);
        console.error(error);
      }
    }
  };
  // Toggle expand/collapse
  // const toggleExpand = (name: string) => {
  //   setExpanded((prevState) => ({ ...prevState, [name]: !prevState[name] }));
  // };

  // Fetch data on mount
  useEffect(() => {
    fetchKeysaasInstances();
  }, []);

  if (loading) {
    return <div>Loading KeySaaS instances...</div>;
  }

  if (error) {
    return <div className="error">{error}</div>;
  }

  return (
      <div className="list-container">
        <h2>KeySaaS Instances</h2>
        <ul className="keysaas-list">
          {keySaaSInstances.map((instance: any) => (

<li
  key={instance.metadata.name}
  className={`keysaas-item expanded`}
  onClick={() => goToDetail(instance.metadata.name)} // Add click handler to the entire item
>
  <div className="keysaas-summary">
    <p><strong>Name:</strong> {instance.metadata.name}</p>
    <p><strong>Status:</strong> {instance.status?.status || "Unknown"}</p>
    
    {/* Delete button */}
    <button
      className="delete-button"
      onClick={(e) => {
        e.stopPropagation(); // Prevent triggering the item click
        handleDelete(instance.metadata.name);
      }}
    >
      Delete
    </button>
  </div>
  
  <div className="keysaas-details">
    <p><strong>Pod Name:</strong> {instance.status?.podName}</p>
    <p><strong>Secret Name:</strong> {instance.status?.secretName}</p>
    <p><strong>URL:</strong> {instance.status?.url}</p>
  </div>
</li>


          
          ))}
        </ul>
      </div>
  );
};

export default KeysaasList;
