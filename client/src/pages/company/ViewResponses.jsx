import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { companyAPI } from '../../services/api';

const ViewResponses = () => {
  const { vacancyId } = useParams();
  const [responses, setResponses] = useState([]);
  const [vacancy, setVacancy] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    fetchResponses();
  }, [vacancyId]);

  const fetchResponses = async () => {
    try {
      const res = await companyAPI.getResponses(vacancyId);
      setResponses(res.data.responses || []);
      setVacancy(res.data.vacancy);
    } catch (err) {
      setError('Ошибка загрузки откликов');
    } finally {
      setLoading(false);
    }
  };

  if (loading) return <div className="loading">Загрузка...</div>;

  return (
    <div className="container">
      <h2>Отклики на вакансию: {vacancy?.name}</h2>
      {error && <div className="error-message">{error}</div>}

      {responses.length === 0 ? (
        <div className="empty-state">
          <p>Пока нет откликов на эту вакансию</p>
        </div>
      ) : (
        <div className="responses-list">
          {responses.map((response) => (
            <div key={response.id} className="response-card">
              <div className="response-header">
                <h3>{response.name}</h3>
                <span className="university-badge">{response.university_name}</span>
              </div>
              <div className="response-body">
                <p><strong>Курс:</strong> {response.course}</p>
                <p><strong>Контакты:</strong> {response.contacts}</p>
                <p><strong>Дата подачи:</strong> {new Date(response.created_at).toLocaleDateString('ru-RU')}</p>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default ViewResponses;
