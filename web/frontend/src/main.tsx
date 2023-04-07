import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App'
import './index.css'

console.warn("TODO: update the wasm.main file to be loaded here")

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
)
