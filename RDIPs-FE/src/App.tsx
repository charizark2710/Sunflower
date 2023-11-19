import { Provider } from 'react-redux';
import { BrowserRouter } from 'react-router-dom';
import './App.scss';
import { store } from './redux/store';
import { CommonRoute } from './routes';

function App() {
  return (
    <Provider store={store}>
      <BrowserRouter>
        <CommonRoute />
      </BrowserRouter>
    </Provider>
  );
}

export default App;
