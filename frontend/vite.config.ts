import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  server: {
    host: '0.0.0.0',
    port: 5173,
    allowedHosts: [
      'localhost',
      '127.0.0.1',
      '192.168.0.8',
      '6a15e648db8c.ngrok-free.app', // <-- ваш ngrok-домен
    ],
  },
  plugins: [react()],

})
