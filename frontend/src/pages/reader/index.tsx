import { useEffect, useState, useRef } from "react";
import { useParams } from "react-router-dom";
import { bookClient, bookmarkClient, noteClient } from "../../connect";
import {
  Typography,
  Paper,
  Box,
  Stack,
  IconButton,
  TextField,
  Menu,
  MenuItem,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  LinearProgress,
} from "@mui/material";
import ArrowBackIosIcon from "@mui/icons-material/ArrowBackIos";
import ArrowForwardIosIcon from "@mui/icons-material/ArrowForwardIos";
import { Document, Page, pdfjs } from "react-pdf";
import "react-pdf/dist/Page/AnnotationLayer.css";
import "react-pdf/dist/Page/TextLayer.css";

pdfjs.GlobalWorkerOptions.workerSrc = "/pdf.worker.min.mjs";

export default function ReaderPage() {
  const { bookId } = useParams<{ bookId: string }>();
  const [book, setBook] = useState<any>(null);
  const [numPages, setNumPages] = useState<number>(0);
  const [pageNumber, setPageNumber] = useState<number>(1);
  const [contextMenu, setContextMenu] = useState<{
    mouseX: number;
    mouseY: number;
    text: string;
    page: number;
  } | null>(null);
  const [noteDialog, setNoteDialog] = useState<{
    open: boolean;
    text: string;
    page: number;
  }>({ open: false, text: "", page: 1 });
  const [noteInput, setNoteInput] = useState("");
  const [inputPage, setInputPage] = useState<number | string>("");
  const [notes, setNotes] = useState<any[]>([]);
  const [bookmarks, setBookmarks] = useState<any[]>([]);
  const [showNotes, setShowNotes] = useState(false);

  const [audioUrl, setAudioUrl] = useState<string | null>(null);
  const [isPlaying, setIsPlaying] = useState(false);
  const [autoListening, setAutoListening] = useState(false);
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);

  const goToPrevPage = () => {
    setPageNumber((prev) => Math.max(prev - 1, 1));
    setAutoListening(false);
    clearTimeout(timeoutRef.current as NodeJS.Timeout);
  };

  const goToNextPage = () => {
    setPageNumber((prev) => {
      const nextPage = Math.min(prev + 1, numPages);
      bookClient.updateBookPage({ bookId, page: nextPage });
      return nextPage;
    });
    setAutoListening(false);
    clearTimeout(timeoutRef.current as NodeJS.Timeout);
  };

  useEffect(() => {
    setInputPage(pageNumber);
  }, [pageNumber]);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setInputPage(e.target.value.replace(/[^0-9]/g, ""));
  };

  const checkAudioExists = async (url: string) => {
    try {
      const res = await fetch(url);
      return res.ok;
    } catch {
      return false;
    }
  };

  // Основная функция для автопрослушивания
  const playOrSkipAudio = async (page: number) => {
    if (!book || !book.audioPath) return;
    let audioBase = book.audioPath.replace(/\/full\.mp3$/, "");
    const url = `/api${audioBase}/pages/page_${page}.mp3`;
    const exists = await checkAudioExists(url);
    if (exists) {
      setAudioUrl(url);
      setIsPlaying(true);
    } else {
      setAudioUrl(null);
      setIsPlaying(false);
      if (autoListening && page < numPages) {
        timeoutRef.current = setTimeout(() => {
          setPageNumber(page + 1);
          setInputPage(page + 1);
        }, 2000);
      } else {
        setAutoListening(false);
      }
    }
  };

  // useEffect для автопрослушивания при смене страницы
  useEffect(() => {
    if (autoListening) {
      playOrSkipAudio(pageNumber);
    } else {
      setAudioUrl(null);
      setIsPlaying(false);
    }
    // eslint-disable-next-line
  }, [pageNumber, autoListening, book]);

  // Очищаем таймер при размонтировании или смене страницы
  useEffect(() => {
    return () => {
      if (timeoutRef.current) clearTimeout(timeoutRef.current);
    };
  }, []);

  // Кнопка "слушать" - старт/стоп
  const handleListenPage = () => {
    if (autoListening) {
      setAutoListening(false);
      setAudioUrl(null);
      setIsPlaying(false);
      clearTimeout(timeoutRef.current as NodeJS.Timeout);
    } else {
      setAutoListening(true);
      playOrSkipAudio(pageNumber);
    }
  };

  const handleAudioEnded = () => {
    if (autoListening && pageNumber < numPages) {
      setPageNumber(pageNumber + 1);
      setInputPage(pageNumber + 1);
    } else {
      setIsPlaying(false);
      setAudioUrl(null);
      setAutoListening(false);
    }
  };

  const handleInputBlur = () => {
    let page = Number(inputPage);
    if (isNaN(page) || page < 1) page = 1;
    if (page > numPages) page = numPages;
    setPageNumber(page);
    bookClient.updateBookPage({ bookId, page });
    setAutoListening(false);
    clearTimeout(timeoutRef.current as NodeJS.Timeout);
  };

  const handleInputKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter") {
      (e.target as HTMLInputElement).blur();
    }
  };

  const handleToggleNotes = () => setShowNotes((prev) => !prev);

  useEffect(() => {
    const fetchBook = async () => {
      const res = await bookClient.getBook({ bookId });
      setBook(res.book);
      const notesRes = await noteClient.getNotes({ bookId });
      setNotes(notesRes.notes || []);
      const bookmarksRes = await bookmarkClient.getBookmarks({ bookId });
      setBookmarks(bookmarksRes.bookmarks || []);
      const initialPage =
        res.book && typeof res.book.page === "number" && res.book.page > 0
          ? res.book.page
          : 1;
      setPageNumber(initialPage);
    };
    fetchBook();
  }, [bookId]);

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === "ArrowLeft" || e.key === "a") {
        goToPrevPage();
      }
      if (e.key === "ArrowRight" || e.key === "d") {
        goToNextPage();
      }
    };
    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
    // eslint-disable-next-line
  }, []);

  if (!book) return <Typography>Загрузка...</Typography>;

  const handlePdfContextMenu = (event: React.MouseEvent) => {
    event.preventDefault();
    const selection = window.getSelection();
    const text = selection ? selection.toString() : "";
    if (text) {
      setContextMenu({
        mouseX: event.clientX - 2,
        mouseY: event.clientY - 4,
        text,
        page: pageNumber,
      });
    }
  };

  const handleAddBookmark = async () => {
    await bookmarkClient.addBookmark({
      bookId,
      page: pageNumber,
      note: "",
    });
    const bookmarksRes = await bookmarkClient.getBookmarks({ bookId });
    setBookmarks(bookmarksRes.bookmarks || []);
  };

  return (
    <Box sx={{ display: "flex", gap: 3, mt: 3 }}>
      <Paper sx={{ flex: 7, p: 3, minWidth: 0 }}>
        <Box
          display="flex"
          justifyContent="space-between"
          alignItems="center"
          mb={2}
        >
          <Box>
            <Typography variant="h4" gutterBottom>
              {book.title}
            </Typography>
            <Typography variant="subtitle1" gutterBottom>
              Автор: {book.author}
            </Typography>
            <Typography variant="subtitle2" gutterBottom>
              Дата добавления:{" "}
              {book.createdAt && new Date(book.createdAt).toLocaleString()}
            </Typography>
          </Box>
          <Box>
            <Button onClick={handleToggleNotes} variant="text">
              {showNotes ? "Скрыть меню" : "Открыть меню"}
            </Button>
          </Box>
        </Box>

        <Box
          sx={{
            mt: 4,
            height: 800,
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
            overflow: "hidden",
            background: "#222",
            borderRadius: 2,
          }}
          onContextMenu={handlePdfContextMenu}
        >
          <Document
            file={`/api${
              book.filePath.startsWith("/")
                ? book.filePath
                : "/" + book.filePath
            }`}
            onLoadSuccess={({ numPages }) => {
              setNumPages(numPages);
              setPageNumber(pageNumber > numPages ? numPages : pageNumber);
            }}
            loading="Загрузка PDF..."
            error="Ошибка загрузки PDF"
          >
            <Page pageNumber={pageNumber} width={600} />
          </Document>
        </Box>
        <Stack
          direction="row"
          spacing={2}
          justifyContent="center"
          alignItems="center"
          sx={{ mt: 2 }}
        >
          <IconButton
            onClick={goToPrevPage}
            disabled={pageNumber <= 1}
            aria-label="delete"
            size="large"
          >
            <ArrowBackIosIcon fontSize="inherit" />
          </IconButton>

          <Typography>
            Страница {pageNumber} из {numPages}
          </Typography>

          <IconButton
            onClick={goToNextPage}
            disabled={pageNumber >= numPages}
            aria-label="delete"
            size="large"
          >
            <ArrowForwardIosIcon fontSize="inherit" />
          </IconButton>
        </Stack>
        <Box sx={{ width: "100%", mt: 2 }}>
          <LinearProgress
            variant="determinate"
            value={numPages > 0 ? (pageNumber / numPages) * 100 : 0}
            sx={{ height: 8, borderRadius: 4, width: "50%", margin: "0 auto" }}
          />
          <Typography
            variant="caption"
            display="block"
            align="center"
            sx={{ mt: 0.5 }}
          >
            Прогресс:{" "}
            {numPages > 0 ? Math.round((pageNumber / numPages) * 100) : 0}%
          </Typography>
        </Box>
        <Stack direction="row" justifyContent="center" sx={{ mt: 2 }}>
          <TextField
            fullWidth
            label="Перейти к странице"
            value={inputPage}
            onChange={handleInputChange}
            onBlur={handleInputBlur}
            onKeyDown={handleInputKeyDown}
            inputProps={{
              min: 1,
              max: numPages,
              style: { width: 80 },
              inputMode: "numeric",
              pattern: "[0-9]*",
            }}
            sx={{
              width: 200,
              "& input[type=number]::-webkit-outer-spin-button, & input[type=number]::-webkit-inner-spin-button":
                {
                  WebkitAppearance: "none",
                  margin: 0,
                },
              "& input[type=number]": {
                MozAppearance: "textfield",
              },
            }}
          />
          <Button variant="outlined" onClick={handleAddBookmark} sx={{ ml: 2 }}>
            Добавить закладку
          </Button>
        </Stack>
        <Stack direction="row" justifyContent="center" sx={{ mt: 2 }}>
          <IconButton
            onClick={handleListenPage}
            aria-label="listen"
            size="large"
            color={autoListening ? "primary" : "default"}
          >
            <span role="img" aria-label="listen">
              🔊
            </span>
          </IconButton>
          {audioUrl && autoListening && (
            <audio
              src={audioUrl}
              autoPlay
              controls
              style={{ width: "100%", marginTop: 16 }}
              onEnded={handleAudioEnded}
            />
          )}
        </Stack>

        <Menu
          open={!!contextMenu}
          onClose={() => setContextMenu(null)}
          anchorReference="anchorPosition"
          anchorPosition={
            contextMenu
              ? { top: contextMenu.mouseY, left: contextMenu.mouseX }
              : undefined
          }
        >
          <MenuItem
            onClick={async () => {
              await noteClient.addNote({
                bookId,
                page: contextMenu!.page,
                text: contextMenu!.text,
              });
              const notesRes = await noteClient.getNotes({ bookId });
              setNotes(notesRes.notes || []);
              setContextMenu(null);
            }}
          >
            Добавить заметку
          </MenuItem>
        </Menu>

        <Dialog
          open={noteDialog.open}
          onClose={() => setNoteDialog({ ...noteDialog, open: false })}
        >
          <DialogTitle>Добавить заметку</DialogTitle>
          <DialogContent>
            <Typography variant="body2" gutterBottom>
              Страница: {noteDialog.page}
            </Typography>
            <Typography variant="body2" gutterBottom>
              Выделено: {noteDialog.text}
            </Typography>
            <TextField
              label="Текст заметки"
              fullWidth
              multiline
              minRows={2}
              value={noteInput}
              onChange={(e) => setNoteInput(e.target.value)}
              sx={{ mt: 2 }}
            />
          </DialogContent>
          <DialogActions>
            <Button
              onClick={() => setNoteDialog({ ...noteDialog, open: false })}
            >
              Отмена
            </Button>
            <Button
              onClick={async () => {
                await noteClient.addNote({
                  bookId,
                  page: noteDialog.page,
                  text: noteInput || noteDialog.text,
                });
                const notesRes = await noteClient.getNotes({ bookId });
                setNotes(notesRes.notes || []);
                setNoteDialog({ ...noteDialog, open: false });
                setNoteInput("");
                const res = await bookClient.getBook({ bookId });
                setBook(res.book);
              }}
              disabled={!noteInput && !noteDialog.text}
              variant="contained"
            >
              Сохранить
            </Button>
          </DialogActions>
        </Dialog>
      </Paper>
      {/* Right: Notes & Bookmarks List */}
      {showNotes && (
        <Paper
          sx={{ flex: 3, p: 3, minWidth: 0, height: 900, overflow: "auto" }}
        >
          <Typography variant="h5" gutterBottom>
            Заметки
          </Typography>
          {notes.length > 0 ? (
            notes.map((note) => (
              <Paper
                key={note.id}
                onClick={() => {
                  setPageNumber(note.page);
                  setInputPage(note.page);
                  setAutoListening(false);
                  clearTimeout(timeoutRef.current as NodeJS.Timeout);
                }}
                sx={{
                  p: 2,
                  mb: 2,
                  cursor: "pointer",
                  "&:hover": { backgroundColor: "action.hover" },
                }}
              >
                <Typography variant="body1">{note.text}</Typography>
                <Typography variant="body2" color="text.secondary">
                  Страница: {note.page}
                </Typography>
              </Paper>
            ))
          ) : (
            <Typography color="text.secondary">Нет заметок</Typography>
          )}

          <Typography variant="h5" gutterBottom sx={{ mt: 4 }}>
            Закладки
          </Typography>
          {bookmarks.length > 0 ? (
            bookmarks.map((bm) => (
              <Paper
                key={bm.id}
                onClick={() => {
                  setPageNumber(bm.page);
                  setInputPage(bm.page);
                  setAutoListening(false);
                  clearTimeout(timeoutRef.current as NodeJS.Timeout);
                }}
                sx={{
                  p: 2,
                  mb: 2,
                  cursor: "pointer",
                  "&:hover": { backgroundColor: "action.hover" },
                }}
              >
                <Typography variant="body2" color="text.secondary">
                  Страница: {bm.page}
                </Typography>
                {bm.note && <Typography variant="body1">{bm.note}</Typography>}
              </Paper>
            ))
          ) : (
            <Typography color="text.secondary">Нет закладок</Typography>
          )}
        </Paper>
      )}
    </Box>
  );
}