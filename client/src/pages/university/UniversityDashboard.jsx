import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { universityAPI } from '../../services/api';

const UniversityDashboard = () => {
  const [applications, setApplications] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    fetchApplications();
  }, []);

  const fetchApplications = async () => {
    try {
      const response = await universityAPI.getMyApplications();
      setApplications(response.data);
    } catch (err) {
      setError('Ошибка загрузки заявок');
    } finally {
      setLoading(false);
    }
  };

  if (loading) return <div className="loading">Загрузка...</div>;

  return (
    <div className="container">
      <div className="dashboard-header">
        <h2>Мои заявки на стажировки</h2>
        <Link to="/university/create-application" className="btn btn-primary">
          + Создать новую заявку
        </Link>
      </div>

      {error && <div className="error-message">{error}</div>}

      {applications.length === 0 ? (
        <div className="empty-state">
          <p>У вас пока нет заявок</p>
          <Link to="/university/create-application" className="btn btn-secondary">
            Создать первую заявку
          </Link>
        </div>
      ) : (
        <div className="applications-grid">
          {applications.map((app) => (
            <div key={app.id} className="application-card">
              <div className="application-header">
                <h3>{app.vacancy_name}</h3>
                <span className={`status-badge status-${app.status}`}>
                  {app.status === 'pending' ? 'На рассмотрении' : 
                   app.status === 'approved' ? 'Одобрено' : 'Отклонено'}
                </span>
              </div>
              <div className="application-body">
                <p><strong>Студент:</strong> {app.name}</p>
                <p><strong>Курс:</strong> {app.course}</p>
                <p><strong>Контакты:</strong> {app.contacts}</p>
                <p><strong>Дата подачи:</strong> {new Date(app.created_at).toLocaleDateString('ru-RU')}</p>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default UniversityDashboard;
