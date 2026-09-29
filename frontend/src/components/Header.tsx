import { Link, useNavigate } from "react-router-dom";
import { isAuthenticated, removeToken } from "../api/client";

export function Header() {
  const navigate = useNavigate();
  const authenticated = isAuthenticated();

  function handleLogout() {
    removeToken();
    navigate("/login");
  }

  return (
    <header className="header">
      <div className="header-inner">
        <Link to="/restaurants" className="logo">
          BiteDash
        </Link>
        <nav className="nav">
          <Link to="/restaurants">Рестораны</Link>
          <Link to="/cart">Корзина</Link>
          <Link to="/orders">Заказы</Link>
          {authenticated ? (
            <button type="button" className="btn btn-link" onClick={handleLogout}>
              Выйти
            </button>
          ) : (
            <>
              <Link to="/login">Вход</Link>
              <Link to="/register">Регистрация</Link>
            </>
          )}
        </nav>
      </div>
    </header>
  );
}
