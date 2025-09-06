import { useEffect, useState } from "react";
import {
  Box,
  Typography,
  Button,
  Paper,
  List,
  ListItem,
  ListItemText,
  Avatar,
} from "@mui/material";
// Import booksClient from its module
import { bookClient } from "../../connect";

const AudioBooksPage = () => {
  const [books, setBooks] = useState<any[]>([]);

  useEffect(() => {
    const fetchBooks = async () => {
      try {
        const res = await bookClient.getBooks({});
        setBooks(res.books || []);
      } catch (error) {
        console.error("Error fetching books:", error);
      }
    };
    fetchBooks();
  }, []);

  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Аудиокниги
      </Typography>
      <List>
        {books.map((book) => (
          <Paper key={book.id} sx={{ mb: 2, p: 2 }}>
            <ListItem>
              {book.coverUrl && (
                <Avatar
                  variant="square"
                  src={`http://localhost:50051${book.coverUrl}`}
                  alt={book.title}
                  sx={{ width: 64, height: 90, mr: 2 }}
                />
              )}
              <ListItemText primary={book.title} secondary={book.author} />
            </ListItem>
            {book.audioPath && (
              <Box sx={{ mt: 1 }}>
                <audio
                  controls
                  src={`http://localhost:50051${book.audioPath}`}
                  style={{ width: "100%" }}
                />

              </Box>
            )}
          </Paper>
        ))}
      </List>
    </Box>
  );
};

export default AudioBooksPage;
