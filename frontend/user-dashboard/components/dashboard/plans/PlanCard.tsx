"use client";
import Image from "next/image";
import { FaDumbbell, FaAppleAlt, FaHeartbeat } from "react-icons/fa";
import type { Plan } from "@/lib/types/plans";
import { getPlanImage, safeUpper } from "./utils";
import { motion } from "framer-motion";


export default function PlanCard({ plan, onSelect }: { plan: Plan; onSelect?: (p: Plan) => void }) {
const icon = (() => {
switch ((plan?.plan_type ?? '').toLowerCase()) {
case "fitness":
case "workout":
return <FaDumbbell className="text-red-500" />;
case "diet":
return <FaAppleAlt className="text-green-500" />;
case "physio":
return <FaHeartbeat className="text-blue-500" />;
default:
return <FaDumbbell className="text-gray-400" />;
}
})();


return (
<motion.button
whileHover={{ scale: 1.03 }}
whileTap={{ scale: 0.98 }}
onClick={() => onSelect?.(plan)}
className="text-left bg-white rounded-2xl shadow hover:shadow-xl transition overflow-hidden focus:outline-none focus:ring-2 focus:ring-blue-500"
>
<div className="relative w-full h-48">
<Image src={getPlanImage(plan?.plan_type)} alt={plan?.name ?? "Plan"} fill className="object-cover" />
</div>
<div className="p-4">
<div className="flex items-center gap-2 text-xl mb-2">
{icon}
<h3 className="text-lg font-semibold line-clamp-1">{plan?.name ?? "Unnamed Plan"}</h3>
</div>
<p className="text-gray-700 text-sm line-clamp-2 min-h-[40px]">{plan?.description ?? ""}</p>
<div className="mt-3 flex items-center gap-2 text-xs">
<span className="px-2 py-1 bg-blue-50 text-blue-700 rounded">{safeUpper(plan?.plan_type)}</span>
<span className="px-2 py-1 bg-gray-50 text-gray-700 rounded">{safeUpper(plan?.goal_type)}</span>
<span className="px-2 py-1 bg-emerald-50 text-emerald-700 rounded">{plan?.duration ?? 0} days</span>
</div>
</div>
</motion.button>
);
}