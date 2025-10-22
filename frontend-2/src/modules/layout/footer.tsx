import { Link } from "react-router-dom"

const NAV = [{
    path: '/',
    label: 'Создать путь'
},
{
    path: '/status',
    label: 'Статусы путей'
}]

export const Footer = () => {
    return (
        <footer className="bg-gray-800 text-white py-4 mt-8">
            <nav className="container mx-auto">
                <ul className="flex justify-center space-x-8">
                    {NAV.map((element) => (
                        <li key={element.label}>
                            <Link 
                                to={element.path} 
                                className="hover:text-blue-400 transition-colors duration-200 text-lg font-medium"
                            >
                                {element.label}
                            </Link>
                        </li>
                    ))}
                </ul>
            </nav>
        </footer>
    )
}
