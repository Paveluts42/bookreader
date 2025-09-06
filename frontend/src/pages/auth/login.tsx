import { useState } from "react";
import { TextField, Button, Paper, Typography } from "@mui/material";
import { useNavigate } from "react-router-dom";
import { userClient } from "../../connect";
import { useAuthStore } from "../../store/auth";

export default function LoginPage() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const navigate = useNavigate();

  const { setTokens,setUser } = useAuthStore();

  const handleLogin = async () => {
    try {
      const res = await userClient.login({ username, password });

      console.log("Login response:", res);
      if (res.userId) {
                setTokens(res.accessToken, res.refreshToken);


        localStorage.setItem("token", res.accessToken);
        localStorage.setItem("refreshToken", res.refreshToken);
        localStorage.setItem("user", res.userId);
       const userRes= await userClient.getUser( { userId: res.userId } ); 
       setUser(userRes);
        navigate("/");
      } else {
        setError(res.error || "Ошибка входа");
      }
    } catch {
      setError("Ошибка входа");
    }
  };
const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    handleLogin();
  };
  return (
    <Paper sx={{ p: 4, maxWidth: 400, mx: "auto", mt: 20 }}>
      <Typography variant="h5" gutterBottom>Вход</Typography>
            <form onSubmit={handleSubmit}>

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
           onKeyDown={e => {
          if (e.key === "Enter") handleLogin();
        }}
        onChange={e => setPassword(e.target.value)}
      />
      {error && <Typography color="error">{error}</Typography>}
      <Button variant="contained" fullWidth sx={{ mt: 2 }} onClick={handleLogin}>
        Войти
      </Button>
      <Button fullWidth sx={{ mt: 1 }} onClick={() => navigate("/register")}>
        Регистрация
      </Button>
      </form>
    </Paper>
  );
}