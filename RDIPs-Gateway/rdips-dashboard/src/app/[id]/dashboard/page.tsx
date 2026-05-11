import ElectricCapacityChart from '@/components/ElectricCapacityChart';

interface DashboardPageProps {
  params: Promise<{
    id: string;
  }>;
}

export default async function DashboardPage({ params }: DashboardPageProps) {
  const { id } = await params;
  return (
    <div className="min-h-screen bg-gray-50 py-8 px-4 sm:px-6 lg:px-8">
      <div className="max-w-6xl mx-auto">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">Device Dashboard</h1>
          <p className="text-gray-600 mt-2">Device ID: {id}</p>
        </div>
        <ElectricCapacityChart id={id} />
      </div>
    </div>
  );
}
