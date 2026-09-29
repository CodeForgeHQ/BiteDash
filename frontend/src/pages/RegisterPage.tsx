import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { getErrorMessage } from "../utils/errors";
import { register } from "../api/authApi";

export function RegisterPage() {
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
      await register({ email, password });
      navigate("/login");
    } catch (err) {
      setError(getErrorMessage(err, "Не удалось зарегистрироваться"));
    }
  }

  return (
    <div className="page page-narrow">
      <h1>Register</h1>
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
          {loading ? "Регистрация..." : "Зарегистрироваться"}
        </button>
      </form>
      <p className="muted">
        Уже есть аккаунт? <Link to="/login">Вход</Link>
      </p>
    </div>
  );
}
