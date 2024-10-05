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

        const result = data.data.result;
        if (result.length === 0) {
          throw new Error('No data returned from Prometheus');
        }
        //const metric = result[0].metric
        let mergedValues : [number, string][] = []
        result.forEach((res : MetricData) => {
          mergedValues = mergedValues.concat(res.values);
        });
        mergedValues = mergedValues.filter((_, index)=> index % Math.ceil(mergedValues.length / MAX_POINTS) === 0);
         
        const labels = mergedValues.map((value: [number, string]) => 
          new Date(convertUnixTimestampToDate(value[0])).toLocaleTimeString([], { 
            hour: '2-digit', 
            minute: '2-digit', 
            second: '2-digit', 
            hour12: false 
          })
        );
        
        const values = mergedValues.map((value: [number, string]) => parseFloat(value[1]));
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
  const formatNumber = (num: number) => {
    if (Math.abs(num) >= 1.0e9) {
      return (num / 1.0e9).toFixed(1) + "b";
    }
    if (Math.abs(num) >= 1.0e6) {
      return (num / 1.0e6).toFixed(1) + "m";
    }
    if (Math.abs(num) >= 1.0e3) {
      return (num / 1.0e3).toFixed(1) + "k";
    }
    return num.toString();
  };

  const options = {
    scales: {
      y: {
        ticks: {
          callback: (value: number | string) => {
            // Format only if the value is a number
            return typeof value === 'number' ? formatNumber(value) : value;
          },
        },
      },  
    },
    plugins: {
      legend: {
        display: true,
      },
    },
    responsive: true,
    maintainAspectRatio: false,
  }
  return (
          <div>   
            <div style={{ height: '400px', width: '100%' }}>

            <Line data={data} options={options} />
            </div>

          </div>
  );
};

function convertUnixTimestampToDate(unixTimestamp: number): Date {
    return new Date(unixTimestamp*1000);
}



export default PrometheusGraph;
