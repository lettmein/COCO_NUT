import '../App.css'
import CreateShipmentPage from "../pages/CreateShipmentPage.tsx";
import { QueryClientProvider } from '@tanstack/react-query';
import {queryClient} from "./providers.ts";

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <CreateShipmentPage/>
    </QueryClientProvider>
  )
}

export default App
