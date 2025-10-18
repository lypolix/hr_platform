import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { universityAPI, vacancyAPI } from '../../services/api';

const CreateApplication = () => {
  const [vacancies, setVacancies] = useState([]);
  const [formData, setFormData] = useState({
    vacancy_id: '',
    name: '',
    course: '',
    contacts: '',
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    fetchVacancies();
  }, []);

  const fetchVacancies = async () => {
    try {
      const response = await vacancyAPI.getAll();
      // Фильтруем только активные вакансии
      const activeVacancies = response.data.filter(v => v.status === 'active');
      setVacancies(activeVacancies);
    } catch (err) {
      setError('Ошибка загрузки вакансий');
    }
  };

  const handleChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      await universityAPI.createApplication({
        ...formData,
        vacancy_id: parseInt(formData.vacancy_id),
        course: parseInt(formData.course),
      });
      alert('Заявка успешно создана!');
      navigate('/university/dashboard');
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка создания заявки');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="container">
      <h2>Создать заявку на стажировку</h2>
      {error && <div className="error-message">{error}</div>}

      <form onSubmit={handleSubmit} className="form-card">
        <div className="form-group">
          <label>Выберите вакансию</label>
          <select
            name="vacancy_id"
            value={formData.vacancy_id}
            onChange={handleChange}
            required
          >
            <option value="">-- Выберите вакансию --</option>
            {vacancies.map((vacancy) => (
              <option key={vacancy.id} value={vacancy.id}>
                {vacancy.name} - {vacancy.description}
              </option>
            ))}
          </select>
        </div>

        <div className="form-group">
          <label>ФИО студента</label>
          <input
            type="text"
            name="name"
            value={formData.name}
            onChange={handleChange}
            required
          />
        </div>

        <div className="form-group">
          <label>Курс</label>
          <input
            type="number"
            name="course"
            value={formData.course}
            onChange={handleChange}
            min="1"
            max="6"
            required
          />
        </div>

        <div className="form-group">
          <label>Контакты</label>
          <textarea
            name="contacts"
            value={formData.contacts}
            onChange={handleChange}
            rows="3"
            required
          />
        </div>

        <button type="submit" className="btn btn-primary" disabled={loading}>
          {loading ? 'Создание...' : 'Создать заявку'}
        </button>
      </form>
    </div>
  );
};

export default CreateApplication;
