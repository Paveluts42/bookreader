import { useEffect, useState, useCallback } from "react";
import { bookClient, userClient } from "../../connect";
import {
  Typography,
  Grid,
  Card,
  CardContent,
  CardMedia,
  CardActionArea,
  Menu,
  MenuItem,
} from "@mui/material";
import { useAuthStore } from "../../store/auth";

export default function BooksList() {
  const { user } = useAuthStore();
  const [books, setBooks] = useState<any[]>([]);
  const [users, setUsers] = useState<any[]>([]);
  const [menuAnchor, setMenuAnchor] = useState<null | HTMLElement>(null);
  const [selectedBookId, setSelectedBookId] = useState<string | null>(null);

  // Получить книги и пользователей (для админа)
  const fetchBooks = useCallback(async () => {
    const res = await bookClient.getBooks({});
    setBooks(res.books || []);
  }, []);

  const fetchUsers = useCallback(async () => {
    if (user?.isAdmin) {
      try {
        const res = await userClient.getUsers({});
        setUsers(res.users || []);
      } catch {
        alert("Ошибка загрузки пользователей");
      }
    }
  }, [user]);

  useEffect(() => {
    fetchBooks();
    fetchUsers();
  }, [fetchBooks, fetchUsers]);

  // Контекстное меню
  const handleContextMenu = (
    event: React.MouseEvent<HTMLDivElement, MouseEvent>,
    bookId: string
  ) => {
    event.preventDefault();
    setMenuAnchor(event.currentTarget);
    setSelectedBookId(bookId);
  };
  const handleCloseMenu = () => {
    setMenuAnchor(null);
    setSelectedBookId(null);
  };

  // Удаление книги
  const handleDelete = async () => {
    if (selectedBookId) {
      await bookClient.deleteBook({ bookId: selectedBookId });
      setBooks((prev) => prev.filter((b) => b.id !== selectedBookId));
    }
    handleCloseMenu();
  };

  // Открытие книги
  const handleOpenBook = () => {
    if (selectedBookId) {
      window.location.href = `/reader/${selectedBookId}`;
    }
    handleCloseMenu();
  };

  // Получить имя владельца книги (для админа)
  const getOwnerName = (userId: string) => {
    const owner = users.find((u) => u.id === userId);
    return owner ? owner.username : "";
  };

  return (
    <>
      <Typography variant="h4" gutterBottom>
        Список книг
      </Typography>
      <Grid container spacing={3}>
        {books.map((b) => (
          <Grid size={3} key={b.id}>
            <Card
              sx={{
                height: 420,
                display: "flex",
                flexDirection: "column",
                boxShadow: 3,
                borderRadius: 3,
                transition: "transform 0.2s",
                "&:hover": { transform: "scale(1.03)", boxShadow: 6 },
              }}
              onContextMenu={(e) => handleContextMenu(e, b.id)}
            >
              <CardActionArea
                sx={{
                  flex: 1,
                  display: "flex",
                  flexDirection: "column",
                  alignItems: "stretch",
                }}
                onClick={() => (window.location.href = `/reader/${b.id}`)}
              >
                <CardMedia
                  component="img"
                  height="220"
                  image={
                    b.coverUrl
                      ? `http://localhost:50051${
                          b.coverUrl.startsWith("/")
                            ? b.coverUrl
                            : "/" + b.coverUrl
                        }`
                      : "/default-book-cover.png"
                  }
                  alt={b.title}
                  sx={{
                    objectFit: "cover",
                    borderTopLeftRadius: 12,
                    borderTopRightRadius: 12,
                  }}
                />
                <CardContent sx={{ flex: 1 }}>
                  <Typography variant="h4" component="div" gutterBottom>
                    {b.title}
                  </Typography>
                  <Typography variant="h6" color="text.secondary" gutterBottom>
                    Автор: {b.author}
                  </Typography>
                  {user?.isAdmin && b.userId && (
                    <Typography variant="body2" color="primary" gutterBottom>
                      Владелец: {getOwnerName(b.userId)}
                    </Typography>
                  )}
                  <Typography variant="caption" color="text.secondary">
                    {typeof b.createdAt === "string"
                      ? `Дата добавления: ${b.createdAt}`
                      : "Дата добавления: —"}
                  </Typography>
                  <br />
                  <Typography variant="caption" color="text.secondary">
                    {typeof b.page === "number" && typeof b.pageAll === "number"
                      ? `Страница: ${b.page} / ${b.pageAll}`
                      : "Страницы: —"}
                  </Typography>
                </CardContent>
              </CardActionArea>
            </Card>
          </Grid>
        ))}
      </Grid>
      <Menu
        open={Boolean(menuAnchor)}
        anchorEl={menuAnchor}
        onClose={handleCloseMenu}
        anchorOrigin={{ vertical: "top", horizontal: "left" }}
        transformOrigin={{ vertical: "top", horizontal: "left" }}
      >
        <MenuItem onClick={handleDelete}>Удалить книгу</MenuItem>
        <MenuItem onClick={handleOpenBook}>Читать книгу</MenuItem>
      </Menu>
    </>
  );
}