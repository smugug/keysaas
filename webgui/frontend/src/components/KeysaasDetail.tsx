import React from 'react';
import './PrometheusGraph';
import PrometheusGraph from './PrometheusGraph';
import { useParams } from 'react-router-dom';


// interface Graph {
//   id: string; // Unique identifier for each graph
//   title: string; // Title for the graph
//   query: string; // The Prometheus query for the graph
// }

const KeysaasDetail: React.FC = () => {
  const { name } = useParams();
  const graphs = [
    { id: '1', title: 'CPU Usage', query: `base_cpu_processCpuTime{job="${name}"}[24h]` },
    { id: '2', title: 'Heap Memory Usage', query: `base_memory_usedHeap_bytes{job="${name}"}[24h]` },
    { id: '3', title: 'GC Pause Time', query: `jvm_gc_pause_seconds_sum{job="${name}"}[24h]` },
    { id: '4', title: 'Memory Usage', query: `vendor_local_container_stats_memory_total{cache_manager="keycloak",job="${name}"}[24h]` },
    // Add more graph configurations as needed  
  ];

  return (
    <div style={{ padding: '20px' }}>
      <h1>Dashboard</h1>
      <div className="flexcontainer">
        {graphs.map((graph) => (
          <div className="flexitem" key={graph.id}>
            <PrometheusGraph query={graph.query} title={graph.title}/>
          </div>
        ))}
      </div>
    </div>
  );
  // return (
  //   <div>
  //     <h1>CUNNY</h1>
  //     <PrometheusGraph 
  //       query='vendor_local_container_stats_memory_total{cache_manager="keycloak"}[24h]'
  //       prometheusUrl="http://kubernetes.local:30003" 
  //     />
  //   </div>
      
    // );
};

export default KeysaasDetail;
