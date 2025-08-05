import { useEffect, useState } from "react";
import { client } from "../../connect";


export default function BooksList() {
  const [books, setBooks] = useState<{ id: string; title: string; author: string }[]>([]);

  useEffect(() => {
    const fetchBooks = async () => {
      const res = await client.getBooks({});
      setBooks(res.books);
    };
    fetchBooks();
  }, []);

  return (
    <div>
      <h2>Список книг</h2>
      <ul>
        {books.map((b) => (
          <li key={b.id}>
            {b.title} — <i>{b.author}</i>
          </li>
        ))}
      </ul>
    </div>
  );
}