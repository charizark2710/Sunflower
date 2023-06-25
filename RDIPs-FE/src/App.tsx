import { BrowserRouter } from 'react-router-dom';
import './App.scss';
import AdminPage from './components/pages/admin/AdminPage';
import { AdminRoute } from './routes';
import { Provider } from 'react-redux';
import { store } from './redux/store';

function App() {
  return (
  <Provider store={store}>
    <BrowserRouter>
      <AdminPage children= {<AdminRoute />} />
    </BrowserRouter>
  </Provider>
  );
}

export default App;
