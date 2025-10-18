import axios from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8081';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Interceptor для добавления токена
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// University API
export const universityAPI = {
  register: (data) => api.post('/universities/sign-up', data),
  login: (data) => api.post('/universities/sign-in', data),
  getProfile: () => api.get('/universities/profile'),
  createApplication: (data) => api.post('/responses', data),
  getMyApplications: () => api.get('/universities/responses'),
};

// Company API
export const companyAPI = {
  register: (data) => api.post('/companies/sign-up', data),
  login: (data) => api.post('/companies/sign-in', data),
  getProfile: () => api.get('/companies/profile'),
  createVacancy: (data) => api.post('/vacancies', data),
  getMyVacancies: () => api.get('/companies/vacancies'),
  getResponses: (vacancyId) => api.get(`/vacancies/${vacancyId}/responses`),
  getAllResponses: () => api.get('/companies/responses'),
};

// Admin API
export const adminAPI = {
  login: (data) => api.post('/admin/sign-in', data),
  getPendingUniversities: () => api.get('/admin/universities/pending'),
  getPendingCompanies: () => api.get('/admin/companies/pending'),
  getPendingVacancies: () => api.get('/admin/vacancies/pending'),
  approveUniversity: (id) => api.put(`/admin/universities/${id}/approve`),
  approveCompany: (id) => api.put(`/admin/companies/${id}/approve`),
  approveVacancy: (id) => api.put(`/admin/vacancies/${id}/approve`),
  rejectUniversity: (id) => api.delete(`/admin/universities/${id}`),
  rejectCompany: (id) => api.delete(`/admin/companies/${id}`),
  rejectVacancy: (id) => api.delete(`/admin/vacancies/${id}`),
};

// Vacancy API (public)
export const vacancyAPI = {
  getAll: () => api.get('/vacancies'),
  getById: (id) => api.get(`/vacancies/${id}`),
};

export default api;
