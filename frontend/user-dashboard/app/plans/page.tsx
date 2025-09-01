"use client";

import { useEffect, useState } from "react";
import Sidebar from "@/components/Sidebar";
import PlanCard from "@/components/dashboard/plans/PlanCard";
import PlanFilters from "@/components/dashboard/plans/PlanFiltersComponent";
import PlanDetails from "@/components/dashboard/plans/PlanDetails";
import PlanActions from "@/components/dashboard/plans/PlanActions";
import type { PlanResponse, PlanFilters as PlanFiltersType } from "@/lib/types/plans";
import apiFetch from "@/lib/api";
import { Dialog, DialogContent } from "@/components/ui/dialog";
import { Loader2 } from "lucide-react";

export default function PlansPage() {
  const [plans, setPlans] = useState<PlanResponse[]>([]);
  const [filteredPlans, setFilteredPlans] = useState<PlanResponse[]>([]);
  const [filters, setFilters] = useState<PlanFiltersType>({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedPlan, setSelectedPlan] = useState<PlanResponse | null>(null);

  useEffect(() => {
    const token = document.cookie
      .split("; ")
      .find((r) => r.startsWith("ft_token="))
      ?.split("=")[1];

    if (!token) {
      setError("You must be logged in to view plans.");
      setLoading(false);
      return;
    }

    apiFetch
      .get<{ message: string; plans: PlanResponse[] | null }>("/plans/available", {
        headers: { Authorization: `Bearer ${token}` },
      })
      .then((res) => {
        setPlans(res.data.plans || []);
        setFilteredPlans(res.data.plans || []);
      })
      .catch(() => setError("Failed to fetch plans."))
      .finally(() => setLoading(false));
  }, []);

  const handleFilterChange = (newFilters: PlanFiltersType) => {
  setFilters(newFilters);
  const q = newFilters.q?.toLowerCase() || "";
  const goal = newFilters.goal_type?.trim() || "";
  const type = newFilters.plan_type?.trim() || "";

  const filtered = plans.filter((p) => {
    const matchesQ =
      p.name.toLowerCase().includes(q) ||
      p.description.toLowerCase().includes(q);
    const matchesGoal = goal ? p.goal_type.trim() === goal : true;
    const matchesType = type ? p.plan_type.trim() === type : true;
    return matchesQ && matchesGoal && matchesType;
  });

  setFilteredPlans(filtered);
};


  return (
    <div className="flex min-h-screen bg-gray-50">
      <aside className="fixed left-0 top-0 h-full w-64">
        <Sidebar />
      </aside>

      <main className="flex-1 md:ml-64 p-6">
        <header className="mb-6">
          <h1 className="text-3xl font-bold">Available Plans</h1>
          <p className="text-gray-600">
            Browse personalized fitness, diet, and physiotherapy templates.
          </p>
        </header>

        <div className="mb-6">
          <PlanFilters value={filters} onChange={handleFilterChange} />
        </div>

        {loading ? (
          <div className="flex justify-center items-center h-40 text-gray-500">
            <Loader2 className="h-6 w-6 animate-spin mr-2" />
            Loading plans...
          </div>
        ) : error ? (
          <div className="text-center text-red-500">{error}</div>
        ) : filteredPlans.length === 0 ? (
          <div className="text-center text-gray-500">No plans match your filters.</div>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
            {filteredPlans.map((plan) => (
              <PlanCard
                key={plan.id}
                plan={plan}
                onSelect={() => setSelectedPlan(plan)}
              />
            ))}
          </div>
        )}
      </main>

      <Dialog open={!!selectedPlan} onOpenChange={() => setSelectedPlan(null)}>
    <DialogContent className="max-w-2xl">
  {selectedPlan && (
    <div className="space-y-4">
      <PlanDetails
        plan={selectedPlan}
        open={!!selectedPlan}
        onClose={() => setSelectedPlan(null)}
        onRequestAssign={(planId: number) =>
          console.log("Assign plan:", planId)
        }
      />
      <PlanActions planId={selectedPlan.id} />
    </div>
  )}
</DialogContent>

      </Dialog>
    </div>
  );
}
