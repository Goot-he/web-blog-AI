import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import App from './App.tsx'
import '../styles/index.css'
import '../styles/site.css'
import '../styles/site-footer.css'
import '../pages/home/HomePage.css'
import '../pages/assistant/AssistantPage.css'
import '../pages/diary/DiaryPage.css'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />
  </StrictMode>,
)
