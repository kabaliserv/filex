import {defineConfig, UserConfig} from 'vite'
import vue from '@vitejs/plugin-vue'
import copy from 'rollup-plugin-copy'
import * as path from 'path'

// https://vitejs.dev/config/

export default defineConfig(({mode, command}) => {
  const config: UserConfig = {
    root: "src",
    plugins: [vue()],
    build: {

      outDir: path.resolve(__dirname, 'dist'),
      emptyOutDir: true,
      // generate manifest.json in outDir
      manifest: true,
      rollupOptions: {
        input: '/main.ts',
      }
    },
    server: {
      port: 3001,
      host: '0.0.0.0'
    }
  }

  if (command == "serve") {
    config.plugins.push(
        copy({
          targets: [
            { src: 'public/**/*', dest: 'dist'},
            { src: 'src/assets/**/*', dest: 'dist/assets' }
          ],
          hook: 'buildStart',
          verbose: true
        })
    )
  }

  return config
})
