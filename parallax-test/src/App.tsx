import './App.css';
import React from 'react';
import { ParallaxProvider } from 'react-scroll-parallax';
import BannerOne from './component/bannerone/BannerOne';
import BannerThree from './component/bannerthree/BannerThree';

function App() {
  return (
    <div id='parallax'>
      <ParallaxProvider>
        <BannerOne />
        <BannerThree/>
      </ParallaxProvider>
    </div>
  );
}

export default App;
