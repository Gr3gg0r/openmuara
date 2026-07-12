import { defineConfig } from 'vite';
import preact from '@preact/preset-vite';
import { compression } from 'vite-plugin-compression2';

export default defineConfig({
  base: '/dashboard-assets/',
  plugins: [
    preact(),
    compression({
      algorithm: 'gzip',
      include: /\.(js|css|html|svg)$/,
      threshold: 1024,
    }),
  ],
  build: {
    outDir: '../../internal/ui/dashboard-dist',
    emptyOutDir: true,
    sourcemap: false,
    minify: 'terser',
    terserOptions: {
      compress: {
        drop_console: true,
        drop_debugger: true,
      },
    },
    rollupOptions: {
      output: {
        entryFileNames: 'assets/[name].js',
        chunkFileNames: 'assets/[name].js',
        assetFileNames: 'assets/[name][extname]',
      },
    },
    reportCompressedSize: true,
  },
  server: {
    port: 5173,
    proxy: {
      '^/(?!@vite|@fs|src/|node_modules/|.*\\.js$|.*\\.css$|.*\\.html$)': {
        target: 'http://127.0.0.1:9000',
        changeOrigin: true,
      },
    },
  },
  resolve: {
    alias: {
      '~': '/src',
    },
  },
  test: {
    environment: 'jsdom',
    globals: true,
    setupFiles: './tests/setup.ts',
  },
});
