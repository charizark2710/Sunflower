import React from "react";
import { useParallax } from "react-scroll-parallax";
import sun from "../../images/sun.png";
import sunlightOne from "../../images/sunlight1.png";
import sunlightTwo from "../../images/sunlight2.png";
import sunlightThree from "../../images/sunlight3.png";
import "./BannerOne.css";

const BannerOne = () => {
  const sunContainer = useParallax<HTMLDivElement>({
    scale: [1.5, 0.2, "easeInQuad"],
    translateX: [-5, 5, "easeInQuad"],
  });

  const sunlightOneParallax = useParallax<HTMLDivElement>({
    translateX: [-5, -200, "easeInQuad"],
  });

  const sunlightTwoParallax = useParallax<HTMLDivElement>({
    translateX: [-5, 200, "easeInQuad"],
  });

  const sunlightThreeParallax = useParallax<HTMLDivElement>({
    translateX: [-5, 10, "easeInQuad"],
  });

  return (
    <div className="banner-one">
      <div className="container">
        <div className="sun" ref={sunContainer.ref}>
          <img src={sun} alt="sun" />
        </div>

        <div className="sunlight-two" ref={sunlightTwoParallax.ref}>
          <img src={sunlightTwo} alt="sunlightTwo" />
        </div>

        <div className="sunlight-one" ref={sunlightOneParallax.ref}>
          <img src={sunlightOne} alt="sunlightOne" />
        </div>

        <div className="sunlight-three" ref={sunlightThreeParallax.ref}>
          <img src={sunlightThree} alt="sunlightThree" />
        </div>
      </div>
    </div>
  );
};

export default BannerOne;
