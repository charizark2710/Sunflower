import React from "react";
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
  rotateY: [-70, 90, "easeOut"],
  // scale: [0.3, 7],
  // translateY: ["290%", "0%"],
  // translateX: ["10%", "70%"],
  translateY: ["50%", "0%"],
  children: (
    <img src={solarDeviceImage} alt="solar device" className="solar-device" />
  ),
};
const BannerThree = () => {
  return (
    <ParallaxBanner
      layers={[sun, yellowBackground, solarDevice]}
      className="banner-three"
    />
  );
};

export default BannerThree;
