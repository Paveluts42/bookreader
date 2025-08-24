import { createClient } from "@connectrpc/connect";
import { createGrpcWebTransport } from "@connectrpc/connect-web";
import { BookService } from "../api/book_pb";
import { NoteService } from "../api/note_pb";
import { UserService } from "../api/user_pb";
import { BookmarkService } from "../api/bookmarks_pb";

let isRefreshing = false;

const transport = createGrpcWebTransport({
  baseUrl: "http://127.0.0.1:50051",
  useBinaryFormat: false,
  interceptors: [
    (next) => async (req) => {
      req.header.set("Content-Type", "application/grpc-web+json");
      const accessToken = localStorage.getItem("token");
      if (accessToken) {
        req.header.set("Authorization", `Bearer ${accessToken}`);
      }


      try {
        const res = await next(req);
        return res;
      } catch (err: any) {


        if (
          err.code === "permission_denied" ||
          err.message === "forbidden" ||
          err.message?.includes("invalid token") ||
          err.message?.includes("missing trailer")||
          err.message?.includes("forbidden")
        ) {
          if (isRefreshing) {
            alert("Сессия истекла. Пожалуйста, войдите снова.");
            localStorage.removeItem("token");
            localStorage.removeItem("refreshToken");
            window.location.href = "/login";
            throw err;
          }
          const refreshToken = localStorage.getItem("refreshToken");
          if (refreshToken) {
            try {
              const refreshRes = await userClient.refreshToken({
                refreshToken,
              });

              if (refreshRes.accessToken && refreshRes.refreshToken) {
                localStorage.setItem("token", refreshRes.accessToken);
                localStorage.setItem("refreshToken", refreshRes.refreshToken);
                req.header.set(
                  "Authorization",
                  `Bearer ${refreshRes.accessToken}`
                );
                console.log(
                  "Повторный запрос с новым токеном:",
                  refreshRes.accessToken
                );
                const retryRes = await next(req);
                return retryRes;
              } else {
                alert("Сессия истекла. Пожалуйста, войдите снова.");
                localStorage.removeItem("token");
                localStorage.removeItem("refreshToken");
                window.location.href = "/login";
                throw err;
              }
            } catch (refreshErr) {
              alert("Ошибка обновления токена. Войдите снова.");
              localStorage.removeItem("token");
              localStorage.removeItem("refreshToken");
              window.location.href = "/login";
              throw refreshErr;
            }
          } else {
            alert("Сессия истекла. Пожалуйста, войдите снова.");
            localStorage.removeItem("token");
            localStorage.removeItem("refreshToken");
            window.location.href = "/login";
            throw err;
          }
        } else {
          alert("Ошибка: " + (err?.message || err));
          throw err;
        }
      }
    },
  ],
});

export const bookmarkClient = createClient(BookmarkService, transport);
export const bookClient = createClient(BookService, transport);
export const noteClient = createClient(NoteService, transport);
export const userClient = createClient(UserService, transport);
