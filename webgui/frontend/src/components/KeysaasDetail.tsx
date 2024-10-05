import React, { useEffect, useState } from 'react';
import './PrometheusGraph';
import PrometheusGraph from './PrometheusGraph';
import { useParams } from 'react-router-dom';
import './KeysaasDetail.css'
import { KeysaasInstance } from './KeysaasList';
import { getKeySaaSInstance } from '../api';
import ThemeUpload from './ThemeUpload';


// interface Graph {
//   id: string; // Unique identifier for each graph
//   title: string; // Title for the graph
//   query: string; // The Prometheus query for the graph
// }

const emptyKeysaasInstance: KeysaasInstance = {
  apiVersion: "",
  kind: "",
  metadata: {
    annotations: {},
    creationTimestamp: "",
    generation: 0,
    name: "",
    namespace: "",
    resourceVersion: "",
    uid: "",
  },
  spec: {
    domainName: "",
    keysaasPassword: "",
    keysaasUsername: "",
    requestsCpu: "",
    requestsMemory: "",
    limitsCpu: "",
    limitsMemory: "",
    scalingThreshold: "",
    minInstances: "",
    maxInstances: "",
    postgresUri: "",
    tls: "",
  },
  status: {
    podName: "",
    secretName: "",
    status: "",
    url: "",
  },
};

const KeysaasDetail: React.FC = () => {
  const { name } = useParams();
  const [keySaaSInstance, setKeySaaSInstance] = useState<KeysaasInstance>(emptyKeysaasInstance);
  useEffect(() => {
    const fetchKeysaasInstance = async () => {
      if (name === undefined) return;
      try {
        const data = await getKeySaaSInstance(name);
        console.log(data);
        
        // Assuming data is a JSON string, parse it to convert to KeysaasInstance
        const parsedData: KeysaasInstance = JSON.parse(data);
        
        // Set the parsed data into state
        setKeySaaSInstance(parsedData);
      } catch (error) {
        console.error(error);
      }
    };

    fetchKeysaasInstance();
  }, [name]); // Just the 'name' dependency is enough here

  const graphs = [
    { id: '1', title: 'CPU Usage', query: `base_cpu_processCpuTime{job="${name}"}[24h]` },
    { id: '2', title: 'Heap Memory Usage', query: `base_memory_usedHeap_bytes{job="${name}"}[24h]` },
    { id: '3', title: 'GC Pause Time', query: `jvm_gc_pause_seconds_sum{job="${name}"}[24h]` },
    { id: '4', title: 'Memory Usage', query: `vendor_local_container_stats_memory_total{cache_manager="keycloak",job="${name}"}[24h]` },
    // Add more graph configurations as needed  
  ];
  return (
    <div className="container">
        <div className="form">
        
        <h2>Details</h2>
        <div className="board">
        <ThemeUpload name={keySaaSInstance.metadata.name} />
        <div className="section">
          <h3>API Version & Kind</h3>
          <p><strong>apiVersion:</strong> {keySaaSInstance.apiVersion}</p>
          <p><strong>kind:</strong> {keySaaSInstance.kind}</p>
        </div>

  <div className="section">
    <h3>Metadata</h3>
    <p><strong>name:</strong> {keySaaSInstance.metadata.name}</p>
    <p><strong>namespace:</strong> {keySaaSInstance.metadata.namespace}</p>
    <p><strong>uid:</strong> {keySaaSInstance.metadata.uid}</p>
    <p><strong>resourceVersion:</strong> {keySaaSInstance.metadata.resourceVersion}</p>
    <p><strong>generation:</strong> {keySaaSInstance.metadata.generation}</p>
    <p><strong>creationTimestamp:</strong> {keySaaSInstance.metadata.creationTimestamp}</p>
  </div>

  <div className="section">
    <h3>Spec</h3>
    <p><strong>domainName:</strong> {keySaaSInstance.spec.domainName}</p>
    <p><strong>keysaasUsername:</strong> {keySaaSInstance.spec.keysaasUsername}</p>
    <p><strong>requestsCpu:</strong> {keySaaSInstance.spec.requestsCpu}</p>
    <p><strong>requestsMemory:</strong> {keySaaSInstance.spec.requestsMemory}</p>
    <p><strong>limitsCpu:</strong> {keySaaSInstance.spec.limitsCpu}</p>
    <p><strong>limitsMemory:</strong> {keySaaSInstance.spec.limitsMemory}</p>
    <p><strong>scalingThreshold:</strong> {keySaaSInstance.spec.scalingThreshold}</p>
    <p><strong>minInstances:</strong> {keySaaSInstance.spec.minInstances}</p>
    <p><strong>maxInstances:</strong> {keySaaSInstance.spec.maxInstances}</p>
    <p><strong>postgresUri:</strong> {keySaaSInstance.spec.postgresUri}</p>
    <p><strong>tls:</strong> {keySaaSInstance.spec.tls}</p>
  </div>

  <div className="section">
    <h3>Status</h3>
    <p><strong>podName:</strong> {keySaaSInstance.status.podName}</p>
    <p><strong>secretName:</strong> {keySaaSInstance.status.secretName}</p>
    <p><strong>status:</strong> {keySaaSInstance.status.status}</p>
    <p><strong>url:</strong> {keySaaSInstance.status.url}</p>
  </div>
</div>

        </div>
          <div className="list">
          <h2>Dashboard</h2>
          <div className="grid-container">
            {graphs.map((graph) => (
              <div className="grid-item" key={graph.id}>
                <PrometheusGraph query={graph.query} title={graph.title} />
              </div>
            ))}
          </div>
        </div>
    </div>    
    
  );
};

export default KeysaasDetail;
