import { useState } from "react";
import { TextField, Button, Paper, Typography } from "@mui/material";
import { useNavigate } from "react-router-dom";
import { userClient } from "../../connect";

export default function RegisterPage() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const navigate = useNavigate();

  const handleRegister = async () => {
    try {
      const res = await userClient.register({ username, password });
      if (res.ok) {
        navigate("/login");
      } else {
        setError(res.error || "Ошибка регистрации");
      }
    } catch {
      setError("Ошибка регистрации");
    }
  };

  return (
    <Paper sx={{ p: 4, maxWidth: 400, mx: "auto", mt: 20 }}>
      <Typography variant="h5" gutterBottom>Регистрация</Typography>
      <TextField
        label="Имя пользователя"
        fullWidth
        margin="normal"
        value={username}
        onChange={e => setUsername(e.target.value)}
      />
      <TextField
        label="Пароль"
        type="password"
        fullWidth
        margin="normal"
        value={password}
        onChange={e => setPassword(e.target.value)}
      />
      {error && <Typography color="error">{error}</Typography>}
      <Button variant="contained" fullWidth sx={{ mt: 2 }} onClick={handleRegister}>
        Зарегистрироваться
      </Button>
      <Button fullWidth sx={{ mt: 1 }} onClick={() => navigate("/login")}>
        Уже есть аккаунт? Войти
      </Button>
    </Paper>
  );
}