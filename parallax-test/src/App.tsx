import "./App.css";
import React from "react";
import { ParallaxProvider } from "react-scroll-parallax";
import BannerZero from "./component/bannerzero/BannerZero";
import BannerOne from "./component/bannerone/BannerOne";
import BannerTwo from "./component/bannertwo/BannerTwo";

function App() {
  return (
    <div id="parallax">
      <ParallaxProvider>
        <BannerZero />
        <BannerOne />
        <BannerTwo />
      </ParallaxProvider>
    </div>
  );
}

export default App;
