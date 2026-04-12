import api from './api'
import type { Monitor, CheckResult, DashboardStats, AuthResponse } from './types'

// Auth
export const login = async (email: string, password: string): Promise<AuthResponse> => {
  const res = await api.post('/auth/login', { email, password })
  return res.data.data  // unwrap data field
}

export const register = async (email: string, password: string): Promise<AuthResponse> => {
  const res = await api.post('/auth/register', { email, password })
  return res.data.data
}

// Monitors
export const getMonitors = async (): Promise<Monitor[]> => {
  const res = await api.get('/monitors')
  return res.data.data
}

export const createMonitor = async (data: {
  name: string
  url: string
  interval_secs: number
}): Promise<Monitor> => {
  const res = await api.post('/monitors', data)
  return res.data.data
}

export const deleteMonitor = async (id: number): Promise<void> => {
  await api.delete(`/monitors/${id}`)
}

// Check results
export const getCheckResults = async (monitorId: number): Promise<CheckResult[]> => {
  const res = await api.get(`/monitors/${monitorId}/results`)
  return res.data.data
}

// Stats
export const getDashboardStats = async (): Promise<DashboardStats> => {
  const res = await api.get('/dashboard/stats')
  return res.data.data
}