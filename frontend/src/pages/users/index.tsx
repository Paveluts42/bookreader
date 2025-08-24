import { useEffect, useState } from "react";
import { Box, Button, Typography, Paper, Avatar, Stack } from "@mui/material";
import { useAuthStore } from "../../store/auth";
import { userClient } from "../../connect";

export default function AdminUsersPage() {
  const { user } = useAuthStore();
  const [users, setUsers] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  const fetchUsers = async () => {
    try {
      const res = await userClient.getUsers({});
      setUsers(res.users || []);
    } catch (err) {
      alert("Ошибка загрузки пользователей");
    } finally {
      setLoading(false);
    }
  };

    useEffect(() => {
    fetchUsers();
  }, []);

  const handleDelete = async (userId: string) => {
    if (userId === user?.userId) return;
    if (!window.confirm("Удалить пользователя?")) return;
    try {
      await userClient.deleteUser({ userId });
      setUsers((prev) => prev.filter((u: any) => u.id !== userId));
    } catch {
      alert("Ошибка удаления пользователя");
    }
  };


  if (!user?.isAdmin) return <Typography>Нет доступа</Typography>;

  return (
    <Box sx={{ mt: 4 }}>
      <Typography variant="h5" gutterBottom>Пользователи</Typography>
      {loading ? (
        <Typography>Загрузка...</Typography>
      ) : (
        <Stack spacing={2}>
          {users.map((u: any) => (
            <Paper key={u.id} sx={{ p: 2, display: "flex", alignItems: "center", gap: 2 }}>
              <Avatar>{u.username[0].toUpperCase()}</Avatar>
              <Box sx={{ flexGrow: 1 }}>
                <Typography>{u.username}</Typography>
                {u.isAdmin && (
                  <Typography variant="caption" color="secondary">Администратор</Typography>
                )}
              </Box>
              {u.id !== user.userId && (
                <Button color="error" onClick={() => handleDelete(u.id)}>
                  Удалить
                </Button>
              )}
            </Paper>
          ))}
        </Stack>
      )}
    </Box>
  );
}