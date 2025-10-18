import React, { useEffect, useState } from 'react';
import { adminAPI } from '../../services/api';

const AdminDashboard = () => {
  const [tab, setTab] = useState('universities'); // 'universities', 'companies', 'vacancies'
  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    fetchData();
  }, [tab]);

  const fetchData = async () => {
    setLoading(true);
    setError('');
    try {
      let response;
      if (tab === 'universities') {
        response = await adminAPI.getPendingUniversities();
      } else if (tab === 'companies') {
        response = await adminAPI.getPendingCompanies();
      } else {
        response = await adminAPI.getPendingVacancies();
      }
      setData(response.data);
    } catch (err) {
      setError('Ошибка загрузки данных');
    } finally {
      setLoading(false);
    }
  };

  const handleApprove = async (id) => {
    try {
      if (tab === 'universities') {
        await adminAPI.approveUniversity(id);
      } else if (tab === 'companies') {
        await adminAPI.approveCompany(id);
      } else {
        await adminAPI.approveVacancy(id);
      }
      alert('Одобрено!');
      fetchData();
    } catch (err) {
      alert('Ошибка одобрения');
    }
  };

  const handleReject = async (id) => {
    if (!window.confirm('Вы уверены, что хотите отклонить?')) return;
    try {
      if (tab === 'universities') {
        await adminAPI.rejectUniversity(id);
      } else if (tab === 'companies') {
        await adminAPI.rejectCompany(id);
      } else {
        await adminAPI.rejectVacancy(id);
      }
      alert('Отклонено!');
      fetchData();
    } catch (err) {
      alert('Ошибка отклонения');
    }
  };

  return (
    <div className="container">
      <h2>Панель администратора</h2>

      <div className="admin-tabs">
        <button
          className={tab === 'universities' ? 'active' : ''}
          onClick={() => setTab('universities')}
        >
          ВУЗы
        </button>
        <button
          className={tab === 'companies' ? 'active' : ''}
          onClick={() => setTab('companies')}
        >
          Компании
        </button>
        <button
          className={tab === 'vacancies' ? 'active' : ''}
          onClick={() => setTab('vacancies')}
        >
          Вакансии
        </button>
      </div>

      {loading && <div className="loading">Загрузка...</div>}
      {error && <div className="error-message">{error}</div>}

      {!loading && data.length === 0 && (
        <div className="empty-state">
          <p>Нет заявок на модерацию</p>
        </div>
      )}

      <div className="admin-list">
        {data.map((item) => (
          <div key={item.id} className="admin-card">
            <h3>{item.title || item.name}</h3>
            <p>{item.description || item.inn}</p>
            {item.contacts && <p><strong>Контакты:</strong> {item.contacts}</p>}
            {item.salary && <p><strong>Зарплата:</strong> {item.salary} руб.</p>}
            <div className="admin-actions">
              <button
                className="btn btn-success"
                onClick={() => handleApprove(item.id)}
              >
                Одобрить
              </button>
              <button
                className="btn btn-danger"
                onClick={() => handleReject(item.id)}
              >
                Отклонить
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default AdminDashboard;
