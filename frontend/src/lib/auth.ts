export const setAuth = (token: string) => {
  localStorage.setItem('token', token)
}

export const getToken = () => {
  return localStorage.getItem('token')
}

export const logout = () => {
  localStorage.removeItem('token')
  window.location.href = '/login'
}

export const isAuthenticated = () => {
  const token = localStorage.getItem('token')
  return token !== null && token !== undefined && token !== ''
}