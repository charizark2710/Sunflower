import ReactDOM from 'react-dom/client';
import { Provider } from 'react-redux';
import App from './App';
import './index.css';
import { setupAxios } from "./axios/axiosClient";
import store from './redux/store';

const root = ReactDOM.createRoot(document.getElementById('root') as HTMLElement);
setupAxios(store);
root.render(
  <Provider store={store}>
    <App />
  </Provider>
);
