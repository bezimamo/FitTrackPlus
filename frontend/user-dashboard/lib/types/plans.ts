export type PlanResponse = {
  id: number;
  name: string;
  description: string;
  goal_type: string;
  plan_type: string;
  duration: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
};

export type PlanFilters = {
  q?: string;
  goal_type?: string;
  plan_type?: string;
};
