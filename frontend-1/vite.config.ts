import path from 'path'
import { defineConfig, loadEnv } from 'vite'

import tailwindcss from '@tailwindcss/vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
    const env = loadEnv(mode, process.cwd())

    return {
        plugins: [react(), tailwindcss()],
        define: {
            global: 'window',
        },
        server: {
            open: true,
            port: 5173,
            proxy: {
                '/api': {
                    target: env.VITE_API_URL,
                    changeOrigin: true,
                    secure: true,
                },
            },
        },
        resolve: {
            alias: {
                '@': path.resolve(__dirname, './src'),
                '@pages': path.resolve(__dirname, './src/pages'),
                '@shared': path.resolve(__dirname, './src/shared'),
            },
        },
    }
})