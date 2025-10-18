import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

const Navbar = () => {
  const { user, userType, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/');
  };

  return (
    <nav className="navbar">
      <div className="navbar-brand">
        <Link to="/">HR Platform Mosprom</Link>
      </div>
      <div className="navbar-menu">
        {!user ? (
          <>
            <Link to="/university/login">Вход для ВУЗа</Link>
            <Link to="/company/login">Вход для Компании</Link>
            <Link to="/admin/login">Вход для Администратора</Link>
          </>
        ) : (
          <>
            {userType === 'university' && (
              <>
                <Link to="/university/dashboard">Мои заявки</Link>
                <Link to="/university/create-application">Создать заявку</Link>
              </>
            )}
            {userType === 'company' && (
              <>
                <Link to="/company/dashboard">Мои вакансии</Link>
                <Link to="/company/create-vacancy">Создать вакансию</Link>
                <Link to="/company/responses">Отклики</Link>
              </>
            )}
            {userType === 'admin' && (
              <Link to="/admin/dashboard">Панель администратора</Link>
            )}
            <button onClick={handleLogout} className="logout-btn">Выйти</button>
          </>
        )}
      </div>
    </nav>
  );
};

export default Navbar;
