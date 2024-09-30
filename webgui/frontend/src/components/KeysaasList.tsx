import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import './KeysaasList.css';
import {deleteKeySaaSInstance,getKeySaaSInstances} from '../api'

// Interface for Keysaas Instance
export interface KeysaasInstance {
  metadata: {
    name: string;
  };
  status?: {
    status: string;
    url: string;
    secretName: string,
    podName: string,
  };
  spec: {
    mySQLServiceName: string;
    mySQLUserName: string;
    mySQLUserPassword: string;
    keysaasAdminEmail: string;
    pvcVolumeName: string;
    domainName: string;
    tls: string;
  };
}

const KeysaasList: React.FC = () => {
  const [keySaaSInstances, setKeySaaSInstances] = useState<KeysaasInstance[]>([]);
  const [expanded, setExpanded] = useState<{ [key: string]: boolean }>({});
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
  const toggleExpand = (name: string) => {
    setExpanded((prevState) => ({ ...prevState, [name]: !prevState[name] }));
  };

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
              className={`keysaas-item ${expanded[instance.metadata.name] ? "expanded" : ""}`}
              onClick={() => toggleExpand(instance.metadata.name)}
            >
              <div className="keysaas-summary">
                <p><strong>Name:</strong> {instance.metadata.name}</p>
                <p><strong>Status:</strong> {instance.status?.status || "Unknown"}</p>
                <button className="delete-button" onClick={()=>goToDetail(instance.metadata.name)} >
                  Go to details
                </button>
                <button
                  className="delete-button"
                  onClick={(e) => {
                    e.stopPropagation();
                    handleDelete(instance.metadata.name);
                  }}
                >
                  Delete
                </button>
              </div>

              {expanded[instance.metadata.name] && (
                <div className="keysaas-details">
                  <p><strong>Pod Name:</strong> {instance.status?.podName}</p>
                  <p><strong>Secret Name:</strong> {instance.status?.secretName}</p>
                  <p><strong>URL:</strong> {instance.status?.url}</p>
                  <p><strong>TLS:</strong> {instance.spec.tls}</p>
                </div>
              )}
            </li>
          ))}
        </ul>
      </div>
  );
};

export default KeysaasList;
