import './App.css';
import React from 'react';
import { ParallaxProvider } from 'react-scroll-parallax';
import BannerZero from './component/bannerzero/BannerZero';

function App() {
  return (
    <div id='parallax'>
      <ParallaxProvider>
        <BannerZero />
      </ParallaxProvider>
    </div>
  );
}

export default App;
