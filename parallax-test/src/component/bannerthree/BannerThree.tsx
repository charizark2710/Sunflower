import React, { useEffect, useRef, useState } from "react";
import {
  Parallax,
  ParallaxBanner,
  ParallaxBannerLayer,
} from "react-scroll-parallax";
import { BannerLayer } from "react-scroll-parallax/dist/components/ParallaxBanner/types";
import "./BannerThree.css";
import sunImage from "../../assets/sunset-thumbnail_transparent.png";
import solarDeviceImage from "../../assets/solar.png";
// import solarDeviceImage from '../../assets/123hello.png';
import yellowBackgroundImage from "../../assets/background-yellow_transparent.png";

const sun: BannerLayer = {
  translateX: ["120%", "0%"],
  scale: [1, 0.7],
  children: <img src={sunImage} alt="sun" />,
};

const yellowBackground: BannerLayer = {
  translateY: ["65%", "-35%"],
  speed: 5,
  startScroll: 30,
  children: <img src={yellowBackgroundImage} className="yellow-background" />,
};

const solarDevice: BannerLayer = {
  rotateY: [-90, 30, "easeOut"],
  // scale: [0.3, 7],
  // translateY: ["290%", "0%"],
  // translateX: ["10%", "70%"],
  translateY: ["50%", "0%"],
  children: (
    <img src={solarDeviceImage} alt="solar device" className="solar-device" />
  ),
  // startScroll: 0,
  // endScroll: 950,
  onChange: (e) => {
    const posY = e.el.children[0].getBoundingClientRect().y;
    console.log("W", window.scrollY);
    console.log("Y", posY);
    if (950 < window.scrollY) {
      e.el.children[0].setAttribute("style", `top: ${window.scrollY - 950}px`)
    } else {
      e.el.children[0].removeAttribute("style")
    }
  }
};



const BannerThree = () => {

  // const solarDevice: BannerLayer = {
  //   rotateY: [-70, 0, "easeOut"],
  //   // scale: [0.3, 7],
  //   // translateY: ["290%", "0%"],
  //   // translateX: ["10%", "70%"],
  //   translateY: ["50%", "0%"],
  //   children: (
  //     <img src={solarDeviceImage} alt="solar device" className="solar-device" />
  //   ),
  //   // startScroll: 100,
  //   // endScroll: 500,
  // };

  return (
    <>
      <ParallaxBanner
        layers={[yellowBackground, sun, solarDevice]}
        className="banner-three"
      />
    </>
  );
};

export default BannerThree;
