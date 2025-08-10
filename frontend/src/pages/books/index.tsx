import { useEffect, useState } from "react";
import { client } from "../../connect";
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

export default function BooksList() {
  const [books, setBooks] = useState<
    { id: string; title: string; author: string; coverUrl?: string; page?: number; pageAll?: number }[]
  >([]);
  const [menuAnchor, setMenuAnchor] = useState<null | HTMLElement>(null);
  const [selectedBookId, setSelectedBookId] = useState<string | null>(null);
  const fetchBooks = async () => {
    const res = await client.getBooks({});
    setBooks(res.books);
  };
  useEffect(() => {
    fetchBooks();
  }, []);



  
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
  const handleDelete = async () => {
    if (selectedBookId) {
      await client.deleteBook({ bookId: selectedBookId });
      fetchBooks();
      setBooks((prev) => prev.filter((b) => b.id !== selectedBookId));
    }
    handleCloseMenu();
  };

  const handleOpenBook = () => {
    if (selectedBookId) {
      window.location.href = `/reader/${selectedBookId}`;
    }
    handleCloseMenu();
  };
  return (
    <>
      <Typography variant="h4" gutterBottom>
        Список книг
      </Typography>
<Grid container spacing={3}>
  {books.map((b) => (
    <Grid size={4} key={b.id}>
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
          sx={{ flex: 1, display: "flex", flexDirection: "column", alignItems: "stretch" }}
          onClick={() => window.location.href = `/reader/${b.id}`}
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
            sx={{ objectFit: "cover", borderTopLeftRadius: 12, borderTopRightRadius: 12 }}
          />
          <CardContent sx={{ flex: 1 }}>
            <Typography variant="h4" component="div" gutterBottom>
              {b.title}
            </Typography>
            <Typography variant="h6" color="text.secondary" gutterBottom>
              {"Автор: " + b.author}
            </Typography>
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
