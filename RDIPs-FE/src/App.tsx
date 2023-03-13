import { BrowserRouter } from 'react-router-dom';
import './App.css';
import Header from './components/pages/common/Header';
import { AdminRoute, ScrollParallaxRoute } from './routes';

function App() {
  return (
  <>
    {/* <Header /> */}
    <BrowserRouter>
      <AdminRoute />
      <ScrollParallaxRoute/>
    </BrowserRouter>
  </>
  );
}

export default App;
