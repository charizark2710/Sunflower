import React from "react";
import { useParallax } from "react-scroll-parallax";
import sun from "../../images/sun.png";
import "./BannerOne.css";

const BannerOne = () => {
  const sunContainer = useParallax<HTMLDivElement>({
    scale: [1.5, 0.5, "easeInQuad"],
    translateX: [-50, 100, "easeInQuad"],
  });

  return (
    <div className="banner-one">
      <div className="container">
        <div className="sun" ref={sunContainer.ref}>
          <img src={sun} alt="sun" />
        </div>
      </div>
    </div>
  );
};

export default BannerOne;
