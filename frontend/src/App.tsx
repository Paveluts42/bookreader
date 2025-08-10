import { BrowserRouter, Routes, Route, Link } from "react-router-dom";
import { Container, AppBar, Toolbar, Button, CssBaseline } from "@mui/material";
import { ThemeProvider, createTheme } from "@mui/material/styles";
import BooksList from "./pages/books";
import UploadBookForm from "./pages/uploads";
import ReaderPage from "./pages/reader";

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
  return (
    <ThemeProvider theme={darkTheme}>
      <CssBaseline />
      <BrowserRouter>
        <AppBar position="static">
          <Toolbar>
            <Button color="inherit" component={Link} to="/">
              Список книг
            </Button>
            <Button color="inherit" component={Link} to="/upload">
              Загрузить книгу
            </Button>
          </Toolbar>
        </AppBar>
        <Container sx={{ mt: 4 }}>
          <Routes>
            <Route path="/" element={<BooksList />} />
            <Route path="/upload" element={<UploadBookForm />} />
            <Route path="/reader/:bookId" element={<ReaderPage />} />
          </Routes>
        </Container>
      </BrowserRouter>
    </ThemeProvider>
  );
}