import App from '@/pages/app'
import { RoutePage } from '@/pages/routePage'
import { StatusPage } from '@/pages/statusPage'
import {createBrowserRouter} from 'react-router-dom'

export const router = createBrowserRouter([{

    path: '/',
    element: <App />,
    children: [{
        index: true,
        element: <RoutePage />
    }, {
        path: '/status',
        element: <StatusPage />
    }]
}])