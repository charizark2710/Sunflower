import { createTheme, ThemeProvider } from '@mui/material';
import { Provider } from 'react-redux';
import { BrowserRouter } from 'react-router-dom';
import './App.scss';
import store from './redux/store';
import { CommonRoute } from './routes';

function App() {
  const theme = createTheme({
    palette: {
      primary: {
        main: '#25205B',
      },
    },
  });

  return (
    <Provider store={store}>
      <ThemeProvider theme={theme}>
        <BrowserRouter>
          <CommonRoute />
        </BrowserRouter>
      </ThemeProvider>
    </Provider>
  );
}

export default App;
