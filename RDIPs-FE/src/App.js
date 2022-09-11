import React from 'react';
import { Router, Route, Routes, BrowserRouter } from 'react-router-dom';

import './App.scss';
import Homepage from './components/pages/homepage';
function App() {
  return (
      <BrowserRouter>
        <Routes>
          <Route path='/' element={<Homepage/>} />
        </Routes>
      </BrowserRouter>
  );
}

export default App;
