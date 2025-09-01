"use client";

export default function PlanActions({ planId }: { planId: number }) {
  const handleAssign = () => {
    console.log("Assign plan with ID:", planId);
    // You can add API call here to assign the plan
  };

  return (
    <div className="flex justify-end">
      <button
        onClick={handleAssign}
        className="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition"
      >
        Confirm Assignment
      </button>
    </div>
  );
}
