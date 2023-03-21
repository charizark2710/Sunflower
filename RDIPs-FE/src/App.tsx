import { BrowserRouter } from 'react-router-dom';
import './App.css';
import AdminPage from './components/pages/admin/AdminPage';
import { AdminRoute } from './routes';

function App() {
  return (
  <>
    <BrowserRouter>
      <AdminPage children= {<AdminRoute />} />
    </BrowserRouter>
  </>
  );
}

export default App;
