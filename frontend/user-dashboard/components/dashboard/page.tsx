// src/app/dashboard/page.tsx
export default function DashboardPage() {
  return (
    <div>
      <h2 className="text-2xl font-bold mb-4">Welcome to Your Dashboard</h2>
      <div className="grid grid-cols-3 gap-4">
        <div className="bg-white shadow rounded p-4">
          <h3 className="text-lg font-semibold">Steps Today</h3>
          <p className="text-2xl font-bold text-blue-600">8,543</p>
        </div>
        <div className="bg-white shadow rounded p-4">
          <h3 className="text-lg font-semibold">Calories Burned</h3>
          <p className="text-2xl font-bold text-green-600">456 kcal</p>
        </div>
        <div className="bg-white shadow rounded p-4">
          <h3 className="text-lg font-semibold">Workout Sessions</h3>
          <p className="text-2xl font-bold text-purple-600">3</p>
        </div>
      </div>
    </div>
  );
}
