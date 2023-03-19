import React from "react";
import { ParallaxBanner } from "react-scroll-parallax";
import { BannerLayer } from "react-scroll-parallax/dist/components/ParallaxBanner/types";
import sunImage from "../../assets/sunset-thumbnail_transparent.png";
import "./BannerThree.css";
import solarDeviceImage from "../../assets/solar.png";
import sunLightImage from "../../assets/sunLight.png";

const sun: BannerLayer = {
  translateX: ["120%", "0%"],
  scale: [1, 0.7],
  children: <img src={sunImage} alt="sun" />,
};

const sunLight: BannerLayer = {
  translateX: ["100%", "-100%"],
  startScroll: 30,
  children: <img src={sunLightImage} className="sun-light" />,
  className: "yellowB",
};

const solarDevice1: BannerLayer = {
  rotateY: [0, 90, "easeOut"],
  translateY: ["50%", "0%"],
  children: (
    <img src={solarDeviceImage} alt="solar device" className="solar-device" />
  ),
};

const solarDevice2: BannerLayer = {
  rotateY: [0, 70, "easeOut"],
  translateY: ["44%", "0%"],
  children: (
    <img src={solarDeviceImage} alt="solar device" className="solar-device-2" />
  ),
  onChange: (e) => {
    if (1150 < window.scrollY) {
      if (1250 > window.scrollY) {
        e.el.children[0].setAttribute(
          "style",
          `scale: ${(window.scrollY - 1150) / 200 + 1}`
        );
      } else {
        e.el.children[0].setAttribute("style", `scale: 1.5; top: ${(window.scrollY-1250)/2}px;`);

      }
    } else {
      e.el.children[0].removeAttribute("style");
    }
  },
};

const wire: BannerLayer = {
  translateY: ["0%", "100%"],
  // children: <img src={wireImg} className="wire"/>
};

const yellowBackground: BannerLayer = {
  children: <div className="yellow-background"></div>
}

const BannerThree = () => {
  return (
    <ParallaxBanner
      layers={[sunLight, sun, solarDevice1, solarDevice2, wire]}
      className="banner-three"
    />
  );
};

export default BannerThree;
