import { Chart, registerables } from 'chart.js';
import React, { useEffect, useState } from 'react';
import { Line } from 'react-chartjs-2';
import { getPrometheusMetrics } from '../api';
import '../constants';


Chart.register(...registerables);

interface PrometheusGraphProps {
  query: string;
  title: string;
}

interface PrometheusData {
  labels: string[];
  datasets: {
    label: string;
    data: number[];
    fill: boolean;
    backgroundColor: string;
    borderColor: string;
  }[];
}
interface MetricData {
  metric: Record<string, string>;
  values: [number, string][];
}


const MAX_POINTS = 200;
const PrometheusGraph: React.FC<PrometheusGraphProps> = ({ query, title }) => {
  const [data, setData] = useState<PrometheusData>({ labels: [], datasets: [] });
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      setError(null);

      try {
        const data = await getPrometheusMetrics(query)
        console.log(data)

        const result = data.data.result;
        if (result.length === 0) {
          throw new Error('No data returned from Prometheus');
        }
        //const metric = result[0].metric
        let mergedValues : [number, string][] = []
        result.forEach((res : MetricData) => {
          console.log(res)
          mergedValues = mergedValues.concat(res.values);
        });
        console.log(mergedValues.length) 
        mergedValues = mergedValues.filter((_, index)=> index % Math.ceil(mergedValues.length / MAX_POINTS) === 0);
         
        const labels = mergedValues.map((value: [number, string]) => new Date(convertUnixTimestampToDate(value[0])).toLocaleTimeString());
        const values = mergedValues.map((value: [number, string]) => parseFloat(value[1]));
        console.log(labels)
        console.log(values)
        setData({
          labels,
          datasets: [
            {
              label: title,
              data: values,
              fill: false,
              backgroundColor: 'rgba(75,192,192,0.4)',
              borderColor: 'rgba(75,192,192,1)',
            },
          ],
        });
      } catch (error: any) {
        setError(error.message);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [query, title]);

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;
  const options = {
    responsive: true,
    maintainAspectRatio: false,
    // scales: {
    //   y: {
    //     beginAtZero: true, // Start the y-axis at zero
    //     max: 100, // Set your desired maximum value here
    //     title: {
    //       display: true,
    //       text: 'Memory Usage (%)' // Title for the y-axis
    //     }
    //   },
    //   x: {
    //     title: {
    //       display: true,
    //       text: 'Time (s)' // Title for the x-axis
    //     }
    //   }
    // }
  }
  return (
    <div>
      <h3>{title}</h3>
      <div style={{ height: '400px', width: '80%' }}>

      <Line data={data} options={options} />
      </div>  

    </div>
  );
};

function convertUnixTimestampToDate(unixTimestamp: number): Date {
    return new Date(unixTimestamp*1000);
}



export default PrometheusGraph;
