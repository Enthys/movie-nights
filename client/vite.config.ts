import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import {viteStaticCopy} from 'vite-plugin-static-copy'
import path from 'path';

console.log('---->   ', path.resolve(__dirname, './src/assets/images') + '/*')

// https://vitejs.dev/config/
export default defineConfig({
  server: {
    port: 5000,
  },
  plugins: [
    react(),
    viteStaticCopy({
      targets: [
        {
          src: path.resolve(__dirname, './src/assets/images') + '/*',
          dest: './assets/images'
        }
      ]
    })
  ],
})
