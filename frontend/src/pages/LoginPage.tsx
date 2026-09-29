import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { setToken } from "../api/client";
import { getErrorMessage } from "../utils/errors";
import { login } from "../api/authApi";

export function LoginPage() {
  const navigate = useNavigate();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      const response = await login({ email, password });
      setToken(response.token);
      navigate("/restaurants");
    } catch (err) {
      setError(getErrorMessage(err, "Не удалось войти"));
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="page page-narrow">
      <h1>Login</h1>
      <form className="form" onSubmit={handleSubmit}>
        <label className="form-field">
          Email
          <input
            type="email"
            className="input"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </label>
        <label className="form-field">
          Password
          <input
            type="password"
            className="input"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            minLength={8}
          />
        </label>
        {error && <p className="error-text">{error}</p>}
        <button type="submit" className="btn btn-primary" disabled={loading}>
          {loading ? "Вход..." : "Войти"}
        </button>
      </form>
      <p className="muted">
        Нет аккаунта? <Link to="/register">Регистрация</Link>
      </p>
    </div>
  );
}
