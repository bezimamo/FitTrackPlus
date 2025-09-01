"use client";
import PlanCard from "./PlanCard";
import type { Plan } from "@/lib/types/plans";


export default function PlanOverview({
plans,
onSelect,
}: {
plans: Plan[];
onSelect: (p: Plan) => void;
}) {
if (!plans?.length) {
return (
<div className="text-center text-gray-500 italic py-10">No plans available.</div>
);
}


return (
<div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
{plans.map((p) => (
<PlanCard key={p.id} plan={p} onSelect={onSelect} />
))}
</div>
);
}