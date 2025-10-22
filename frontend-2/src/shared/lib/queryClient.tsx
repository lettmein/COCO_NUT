import { QueryClient, QueryClientProvider } from "@tanstack/react-query"
import { useRef, type ReactNode } from "react"

export const QueryProvider = ({children} : {children: ReactNode}) => {
    const client = useRef(new QueryClient())

    return <QueryClientProvider client={client.current}>{children}</QueryClientProvider>
}