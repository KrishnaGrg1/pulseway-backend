export interface APIResponse<T> {
  success: boolean
  message: string
  data?: T
  error?: {
    code: string
    details: string
  }
}

export interface User {
  id: number
  email: string
}

export interface Monitor {
  id: number
  user_id: number
  name: string
  url: string
  interval_secs: number
  is_active: boolean
  created_at: string
}

export interface CheckResult {
  id: number
  monitor_id: number
  status: 'up' | 'down'
  latency_ms: number
  status_code: number | null
  checked_at: string
}

export interface Incident {
  id: number
  monitor_id: number
  started_at: string
  resolved_at: string | null
  notified: boolean
}

export interface DashboardStats {
  total_monitors: number
  active_monitors: number
  uptime_percentage: number
  avg_latency_ms: number
}

export interface AuthResponse {
  token: string
  user: User
}