"use client";
import { useCallback, useEffect, useMemo, useState } from "react";
import apiFetch from "@/lib/api";
import type { Plan } from "@/lib/types/plans";


export type PlanFilters = {
goal_type?: string;
plan_type?: string;
q?: string; // search query
};


export function usePlans(initialFilters?: PlanFilters) {
const [plans, setPlans] = useState<Plan[]>([]);
const [loading, setLoading] = useState(true);
const [error, setError] = useState<string | null>(null);
const [filters, setFilters] = useState<PlanFilters>(initialFilters ?? {});


const params = useMemo(() => {
const p = new URLSearchParams();
if (filters.goal_type) p.set("goal_type", filters.goal_type);
if (filters.plan_type) p.set("plan_type", filters.plan_type);
return p.toString();
}, [filters]);


const fetchPlans = useCallback(async () => {
try {
setLoading(true);
setError(null);


const token = document.cookie
.split("; ")
.find((row) => row.startsWith("ft_token="))
?.split("=")[1];


if (!token) {
setPlans([]);
setError("You are not logged in");
return;
}


const url = `/plans/available${params ? `?${params}` : ""}`;
const res = await apiFetch.get<{ plans: Plan[] | null }>(url, {
headers: { Authorization: `Bearer ${token}` },
});


setPlans(Array.isArray(res.data.plans) ? res.data.plans : []);
} catch (e) {
setError("Failed to load plans");
setPlans([]);
} finally {
setLoading(false);
}
}, [params]);


useEffect(() => {
fetchPlans();
}, [fetchPlans]);


}