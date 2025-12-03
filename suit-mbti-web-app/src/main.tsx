import React from 'react';
import ReactDOM from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import App from './App';

const queryClient = new QueryClient();

// Vercelデプロイ時はbase pathを削除（動的サイト対応）
// GitHub Pagesデプロイ時は '/Tailorcloud' を使用
const basename = import.meta.env.VERCEL ? '' : (import.meta.env.PROD ? '/Tailorcloud' : '');

const rootElement = document.getElementById('root');

if (rootElement) {
  const root = ReactDOM.createRoot(rootElement);

  root.render(
    <React.StrictMode>
      <QueryClientProvider client={queryClient}>
        <BrowserRouter basename={basename}>
          <App />
        </BrowserRouter>
      </QueryClientProvider>
    </React.StrictMode>
  );
}


