import { Outlet } from "react-router-dom";


export default function App() {
 

  return (
    <div className="p-4 space-y-6">
        <Outlet />
    </div>
  );
}
