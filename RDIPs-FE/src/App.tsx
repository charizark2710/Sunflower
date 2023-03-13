import { BrowserRouter } from 'react-router-dom';
import './App.css';
import Header from './components/pages/common/Header';
import { AdminRoute } from './routes';

function App() {
  return (
  <>
    <Header />
    <BrowserRouter>
      <AdminRoute />
    </BrowserRouter>
  </>
  );
}

export default App;
