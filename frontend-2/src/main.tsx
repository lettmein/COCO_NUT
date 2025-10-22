import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'

import './style.css'
import { QueryProvider } from './shared/lib/queryClient';
import { RouterProvider } from 'react-router-dom';
import { router } from './shared/lib/router';

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <QueryProvider>
            <RouterProvider router={router} />
        </QueryProvider>
    </StrictMode>
)
