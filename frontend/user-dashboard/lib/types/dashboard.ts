export interface UserInfo {
  id: number
  first_name: string
  last_name: string
  email: string
  role: "member" | "trainer" | "admin"
  is_active: boolean
}

export interface Stat {
  total_users: number
  active_plans: number
  total_sessions: number
  completion_rate: number
}

export interface Activity {
  id: number
  type: string
  description: string
  created_at: string
}

export interface Action {
  id: string
  title: string
  description: string
  icon: string
  url: string
  color: string
}

export interface AppNotification {
  id: number
  type: "info" | "warning" | "success" | "error"
  title: string
  message: string
  is_read: boolean
  created_at: string
}

export interface Plan {
  id: number
  name: string
  type: string
  status: string
  progress: number
}

export interface ProgressSummary {
  current_weight: number
  target_weight: number
  weight_lost: number
  overall_progress: number
}

export interface MemberData {
  current_plan: Plan
  progress_summary: ProgressSummary
  upcoming_sessions: any[]
  goals: any[]
}

export interface DashboardData {
  user_role: string
  user_info: UserInfo
  stats: Stat
  recent_activity: Activity[]
  quick_actions: Action[]
  notifications: AppNotification[]
  member_data?: MemberData
}
