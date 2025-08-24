import {
  AppBar,
  Avatar,
  Box,
  Button,
  Container,
  Toolbar,
  Typography,
} from "@mui/material";
import { Link, Route, Routes } from "react-router-dom";
import BooksList from "./books";
import UploadBookForm from "./uploads";
import ReaderPage from "./reader";
import { useAuthStore } from "../store/auth";
import { useEffect } from "react";
import { userClient } from "../connect";
import AdminUsersPage from "./users";

const MainContainer = () => {
  const { user,setUser, logout } = useAuthStore();


  const fetchUser = async () => {
    const userId = localStorage.getItem("user");
    if (userId) {
      try {
       const userRes= await userClient.getUser( { userId: userId } ); 
         setUser(userRes);
      } catch (error) {
        console.error("Error fetching user:", error);
      }
    }
  };
    useEffect(() => {
          fetchUser();
    }, []);


  return (
    <>
      <AppBar position="static">
        <Toolbar sx={{ display: "flex", justifyContent: "space-between" }}>
          <Box sx={{ display: "flex", alignItems: "center", gap: 2 }}>
            <Button color="inherit" component={Link} to="/">
              Список книг
            </Button>
            <Button color="inherit" component={Link} to="/upload">
              Загрузить книгу
            </Button>
               {user?.isAdmin && (
              <Button color="inherit" component={Link} to="/admin/users">
                Пользователи
              </Button>
            )}
          </Box>
          {user && (
            <Box sx={{ display: "flex", alignItems: "center", gap: 1.5 }}>
              <Avatar
                sx={{
                  bgcolor: user.isAdmin ? "secondary.main" : "primary.main",
                  width: 36,
                  height: 36,
                  fontSize: 20,
                }}
              >
                {user.username[0].toUpperCase()}
              </Avatar>
              <Box sx={{ ml: 1 }}>
                <Typography variant="subtitle2" sx={{ fontWeight: 500 }}>
                  {user.username}
                </Typography>
                {user.isAdmin && (
                  <Typography
                    variant="caption"
                    color="secondary"
                    sx={{ lineHeight: 1 }}
                  >
                    Администратор
                  </Typography>
                )}
              </Box>
              <Button color="inherit" onClick={logout} sx={{ ml: 2 }}>
                Выйти
              </Button>
            </Box>
          )}
        </Toolbar>
      </AppBar>
      <Container sx={{ mt: 4 }}>
        <Routes>
          <Route path="/" element={<BooksList />} />
          <Route path="/upload" element={<UploadBookForm />} />
          <Route path="/reader/:bookId" element={<ReaderPage />} />
                {user?.isAdmin && (
            <Route path="/admin/users" element={<AdminUsersPage />} />
          )}
        </Routes>
      </Container>
    </>
  );
};
export default MainContainer;
