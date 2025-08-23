import { BrowserRouter, Routes, Route, Link, Navigate } from "react-router-dom";
import { Container, AppBar, Toolbar, Button, CssBaseline } from "@mui/material";
import { ThemeProvider, createTheme } from "@mui/material/styles";
import BooksList from "./pages/books";
import UploadBookForm from "./pages/uploads";
import ReaderPage from "./pages/reader";
import LoginPage from "./pages/auth/login";
import RegisterPage from "./pages/auth/register";
import ProtectedRoute from "./shared/ProtectedRoute";
import { useAuthStore } from "./store/auth";
import MainContainer from "./pages/mainContainer";

const darkTheme = createTheme({
  palette: {
    mode: "dark",
    background: {
      default: "#121212",
      paper: "#1a1a1a",
    },
  },
});

export default function App() {
    const { token } = useAuthStore();

  return (
    <ThemeProvider theme={darkTheme}>
      <CssBaseline />
      <BrowserRouter>
        {token ? (
          <>
            
            <MainContainer />
          </>
        ) : (
          <Container sx={{ mt: 4 }}>
            <Routes>
              <Route path="/login" element={<LoginPage />} />
              <Route path="/register" element={<RegisterPage />} />
              <Route path="*" element={<Navigate to="/login" />} />
            </Routes>
          </Container>
        )}
      </BrowserRouter>
    </ThemeProvider>
  );
}
