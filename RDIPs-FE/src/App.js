import React, { useState } from 'react';
import { Route, Routes, BrowserRouter } from 'react-router-dom';

import './App.scss';


import Header from './components/pages/common/Header';
import TempPage from './components/pages/Temp.page';
function App() {
  return (
    <>
      <Header />
      <BrowserRouter>
        <Routes>
          <Route path='/' element={<TempPage />} />
        </Routes>
      </BrowserRouter>
     
    </>
   
  );
}

export default App;


