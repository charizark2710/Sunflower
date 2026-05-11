'use client';

import { useEffect, useState } from 'react';
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts';

interface PerformanceData {
  [key: string]: string | number;
}

interface ElectricCapacityChartProps {
  id: string;
}

export default function ElectricCapacityChart({ id }: ElectricCapacityChartProps) {
  const [data, setData] = useState<PerformanceData[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [fields, setFields] = useState<string[]>([]);

  useEffect(() => {
    const fetchPerformances = async () => {
      try {
        setLoading(true);
        const response = await fetch(`/api/performances/${id}`);

        if (!response.ok) {
          throw new Error('Failed to fetch data');
        }

        const result = await response.json();
        
        // Extract Payload from response.data.Payload
        const performanceData = result.data?.Payload || [];
        
        if (performanceData.length > 0) {
          // Transform data to include formatted timestamp
          const transformedData = performanceData.map((item: any) => {
            const timestamp = item.timestamp?.$date?.$numberLong
              ? new Date(parseInt(item.timestamp.$date.$numberLong))
              : new Date(item.putAt);
            
            return {
              ...item,
              timestamp: timestamp.toLocaleTimeString(),
              timestampFull: timestamp.toLocaleString(),
              // Parse numeric values
              capacity: parseFloat(item.capacity),
              inCapacity: parseFloat(item.inCapacity),
              maxCapacity: parseFloat(item.maxCapacity),
              outCapacity: parseFloat(item.outCapacity),
            };
          });
          
          // Extract numeric fields only (exclude metadata fields)
          const dataFields = ['capacity', 'inCapacity', 'maxCapacity', 'outCapacity'];
          
          setFields(dataFields);
          setData(transformedData);
        }
      } catch (err) {
        setError(err instanceof Error ? err.message : 'An error occurred');
        console.error('Error fetching performances:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchPerformances();
  }, [id]);

  if (loading) {
    return (
      <div className="flex items-center justify-center h-96">
        <p className="text-gray-500">Loading chart data...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex items-center justify-center h-96">
        <p className="text-red-500">Error: {error}</p>
      </div>
    );
  }

  if (data.length === 0) {
    return (
      <div className="flex items-center justify-center h-96">
        <p className="text-gray-500">No data available</p>
      </div>
    );
  }

  // Color palette for multiple lines
  const colors = [
    '#8884d8',
    '#82ca9d',
    '#ffc658',
    '#ff7c7c',
  ];

  // Field labels for chart
  const fieldLabels: { [key: string]: string } = {
    capacity: 'Capacity (kW)',
    inCapacity: 'In Capacity (kW)',
    maxCapacity: 'Max Capacity (kW)',
    outCapacity: 'Out Capacity (kW)',
  };

  return (
    <div className="w-full bg-white p-6 rounded-lg shadow-lg">
      <h2 className="text-2xl font-bold mb-6 text-gray-800">Electric Capacity</h2>
      <ResponsiveContainer width="100%" height={400}>
        <LineChart data={data} margin={{ top: 5, right: 30, left: 0, bottom: 5 }}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis
            dataKey="timestamp"
            tick={{ fontSize: 12 }}
            angle={-45}
            textAnchor="end"
            height={80}
          />
          <YAxis tick={{ fontSize: 12 }} />
          <Tooltip
            formatter={(value) => {
              if (typeof value === 'number') {
                return value.toFixed(2);
              }
              return value;
            }}
            labelFormatter={(label) => `Time: ${label}`}
            contentStyle={{ backgroundColor: '#fff', border: '1px solid #ccc' }}
          />
          <Legend />
          {fields.map((field, index) => (
            <Line
              key={field}
              type="monotone"
              dataKey={field}
              stroke={colors[index % colors.length]}
              dot={{ r: 3 }}
              activeDot={{ r: 5 }}
              name={fieldLabels[field] || field}
            />
          ))}
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}
