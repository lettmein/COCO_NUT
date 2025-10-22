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
        <nav>
            <ul>
                {NAV.map((element) => (
                    <li key={element.label}><Link to={element.path}>{element.label}</Link></li>

                ))}
            </ul>
        </nav>
    )
}