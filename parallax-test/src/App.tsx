import './App.css';
import React from 'react';
import { ParallaxProvider } from 'react-scroll-parallax';
import BannerOne from './component/bannerone/BannerOne';

function App() {
  return (
    <div id='parallax'>
      <ParallaxProvider>
        <BannerOne />
      </ParallaxProvider>
    </div>
  );
}

export default App;
