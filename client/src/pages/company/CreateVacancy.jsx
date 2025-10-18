import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { companyAPI } from '../../services/api';

const CreateVacancy = () => {
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    salary: '',
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      await companyAPI.createVacancy({
        ...formData,
        salary: parseFloat(formData.salary),
      });
      alert('Вакансия создана! Ожидайте модерации администратора.');
      navigate('/company/dashboard');
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка создания вакансии');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="container">
      <h2>Создать вакансию</h2>
      {error && <div className="error-message">{error}</div>}

      <form onSubmit={handleSubmit} className="form-card">
        <div className="form-group">
          <label>Название вакансии</label>
          <input
            type="text"
            name="name"
            value={formData.name}
            onChange={handleChange}
            placeholder="Например: Стажер-программист"
            required
          />
        </div>

        <div className="form-group">
          <label>Описание</label>
          <textarea
            name="description"
            value={formData.description}
            onChange={handleChange}
            rows="6"
            placeholder="Опишите требования и условия стажировки"
            required
          />
        </div>

        <div className="form-group">
          <label>Зарплата (руб.)</label>
          <input
            type="number"
            name="salary"
            value={formData.salary}
            onChange={handleChange}
            placeholder="Например: 30000"
            min="0"
            step="1000"
            required
          />
        </div>

        <button type="submit" className="btn btn-primary" disabled={loading}>
          {loading ? 'Создание...' : 'Создать вакансию'}
        </button>
      </form>
    </div>
  );
};

export default CreateVacancy;
