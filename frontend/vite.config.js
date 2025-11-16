import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],

  server: {
    host: '0.0.0.0',
    port: 5173,
    // Enable HMR for faster development
    hmr: true
  },

  // Optimize build output
  build: {
    // Enable source maps for debugging (disable in production)
    sourcemap: false,

    // Chunk size warning limit
    chunkSizeWarningLimit: 1000,

    // Manual chunk splitting for better caching
    rollupOptions: {
      output: {
        manualChunks: {
          // Vendor chunks for better caching
          'react-vendor': ['react', 'react-dom', 'react-router-dom'],
          'lucide-icons': ['lucide-react'],
          'web3-vendor': ['ethers'],
        }
      }
    },

    // Minification settings
    minify: 'esbuild',
    target: 'esnext',

    // Asset optimization
    assetsInlineLimit: 4096 // 4kb - inline small assets
  },

  // Optimize dependencies
  optimizeDeps: {
    include: [
      'react',
      'react-dom',
      'react-router-dom',
      'lucide-react'
    ],
    // Force pre-bundling of these deps
    force: false
  },

  // Resolve configuration
  resolve: {
    // Use browser versions of packages when available
    mainFields: ['browser', 'module', 'main'],
    extensions: ['.js', '.jsx', '.json']
  },

  // Performance optimizations
  esbuild: {
    // Remove console.log in production
    drop: process.env.NODE_ENV === 'production' ? ['console', 'debugger'] : [],
    logOverride: { 'this-is-undefined-in-esm': 'silent' }
  }
})
