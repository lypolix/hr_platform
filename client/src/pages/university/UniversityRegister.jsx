import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { universityAPI } from '../../services/api';

const UniversityRegister = () => {
  const [formData, setFormData] = useState({
    login: '',
    password: '',
    title: '',
    inn: '',
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
      await universityAPI.register(formData);
      alert('Регистрация успешна! Ожидайте подтверждения администратора.');
      navigate('/university/login');
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка регистрации');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="auth-container">
      <div className="auth-card">
        <h2>Регистрация ВУЗа</h2>
        {error && <div className="error-message">{error}</div>}
        
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label>Логин</label>
            <input
              type="text"
              name="login"
              value={formData.login}
              onChange={handleChange}
              required
            />
          </div>

          <div className="form-group">
            <label>Пароль</label>
            <input
              type="password"
              name="password"
              value={formData.password}
              onChange={handleChange}
              required
            />
          </div>

          <div className="form-group">
            <label>Название ВУЗа</label>
            <input
              type="text"
              name="title"
              value={formData.title}
              onChange={handleChange}
              required
            />
          </div>

          <div className="form-group">
            <label>ИНН</label>
            <input
              type="text"
              name="inn"
              value={formData.inn}
              onChange={handleChange}
              required
              maxLength={10}
            />
          </div>

          <button type="submit" className="btn btn-primary" disabled={loading}>
            {loading ? 'Регистрация...' : 'Зарегистрироваться'}
          </button>
        </form>

        <p className="auth-footer">
          Уже есть аккаунт? <Link to="/university/login">Войти</Link>
        </p>
      </div>
    </div>
  );
};

export default UniversityRegister;
