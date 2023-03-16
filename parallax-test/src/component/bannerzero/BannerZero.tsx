import React from "react";
import "./BannerZero.css";
import KeyboardDoubleArrowDownIcon from "@mui/icons-material/KeyboardDoubleArrowDown";

const BannerZero = () => {
  return (
    <div className="banner-zero">
      <div className="container">
        <p className="scroll-text">SCROLL</p>
        <div className="arrow-down-container">
          <KeyboardDoubleArrowDownIcon className="arrow-down" />
        </div>
      </div>
    </div>
  );
};

export default BannerZero;
