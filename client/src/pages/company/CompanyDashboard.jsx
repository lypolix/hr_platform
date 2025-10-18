import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { companyAPI } from '../../services/api';

const CompanyDashboard = () => {
  const [vacancies, setVacancies] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [filter, setFilter] = useState('all'); // 'all', 'active', 'inactive'

  useEffect(() => {
    fetchVacancies();
  }, []);

  const fetchVacancies = async () => {
    try {
      const response = await companyAPI.getMyVacancies();
      setVacancies(response.data);
    } catch (err) {
      setError('Ошибка загрузки вакансий');
    } finally {
      setLoading(false);
    }
  };

  const filteredVacancies = vacancies.filter((v) => {
    if (filter === 'active') return v.status === 'active';
    if (filter === 'inactive') return v.status === 'inactive';
    return true;
  });

  if (loading) return <div className="loading">Загрузка...</div>;

  return (
    <div className="container">
      <div className="dashboard-header">
        <h2>Мои вакансии</h2>
        <Link to="/company/create-vacancy" className="btn btn-primary">
          + Создать вакансию
        </Link>
      </div>

      <div className="filter-tabs">
        <button
          className={filter === 'all' ? 'active' : ''}
          onClick={() => setFilter('all')}
        >
          Все
        </button>
        <button
          className={filter === 'active' ? 'active' : ''}
          onClick={() => setFilter('active')}
        >
          Активные
        </button>
        <button
          className={filter === 'inactive' ? 'active' : ''}
          onClick={() => setFilter('inactive')}
        >
          Неактивные
        </button>
      </div>

      {error && <div className="error-message">{error}</div>}

      {filteredVacancies.length === 0 ? (
        <div className="empty-state">
          <p>Нет вакансий</p>
        </div>
      ) : (
        <div className="vacancies-grid">
          {filteredVacancies.map((vacancy) => (
            <div key={vacancy.id} className="vacancy-card">
              <div className="vacancy-header">
                <h3>{vacancy.name}</h3>
                <span className={`status-badge status-${vacancy.status}`}>
                  {vacancy.status === 'active' ? 'Активна' : 'Неактивна'}
                </span>
              </div>
              <p>{vacancy.description}</p>
              <p><strong>Откликов:</strong> {vacancy.responses_count || 0}</p>
              <Link 
                to={`/company/responses/${vacancy.id}`} 
                className="btn btn-secondary"
              >
                Посмотреть отклики
              </Link>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default CompanyDashboard;
