import React, { useState } from "react";
import "./Sidebar.organism.scss";
// import { images } from "../Helpers/CarouselData";
import ArrowBackIosIcon from "@mui/icons-material/ArrowBackIos";
import ArrowForwardIosIcon from "@mui/icons-material/ArrowForwardIos";

const images = [
  {
    img: 'https://demoda.vn/wp-content/uploads/2022/02/background-dep-1.jpg',
    title: 'Cuoc song ma',
    subtitle: 'Khong nhu la mo'
  }
]

function SidebarOrganism() {
  const [currImg, setCurrImg] = useState(0);
  return (
    <div className="Sidebar">
      <div
        className="SidebarInner"
        style={{ backgroundImage: `url(${images[currImg].img})` }}
      >
        <div
          className="left"
          onClick={() => {
            currImg > 0 && setCurrImg(currImg - 1);
          }}
        >
          <ArrowBackIosIcon style={{ fontSize: 30 }} />
        </div>
        <div className="center">
          <h1>{images[currImg].title}</h1>
          <p>{images[currImg].subtitle}</p>
        </div>
        <div
          className="right"
          onClick={() => {
            currImg < images.length - 1 && setCurrImg(currImg + 1);
          }}
        >
          <ArrowForwardIosIcon style={{ fontSize: 30 }} />
        </div>
      </div>
    </div>
  );
}

export default SidebarOrganism;