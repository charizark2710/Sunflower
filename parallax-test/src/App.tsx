import './App.css';
import React from 'react';
import { ParallaxBanner, ParallaxProvider, useParallax } from 'react-scroll-parallax';
import { BannerLayer } from 'react-scroll-parallax/dist/components/ParallaxBanner/types';
import f7YTTK from '../public/f7YTTK.webp'
function App() {
  return (
    <div id='root' style={
      { background: "#417cdc", width: "100wh", height: "200vh" }}>
      <ParallaxProvider>
        <Parallax />
      </ParallaxProvider>
    </div>
  );
}

const Parallax = (props: any) => {
  // const parallax = useParallax<HTMLDivElement>({ speed: 10 });
  const image: BannerLayer = {
    opacity: [1, 0.3],
    translateY: ['0%', '100%'],
    scale: [1.05, 1, 'easeOutCubic'],
    shouldAlwaysCompleteAnimation: true,
    children: (
      <img src={f7YTTK} />
    )
  };
  const text: BannerLayer = {
    shouldAlwaysCompleteAnimation: true,
    translateY: ['0%', '100%'],
    scale: [1, 1.05, 'easeOutCubic'],
    children: (
      <h1> TEST </h1>
    )
  }

  const paragraph: BannerLayer = {
    translateY: ['100%', '0%'],
    children: (
      <p> {
        `Lorem ipsum dolor sit amet. Quo neque voluptatem non assumenda mollitia in nostrum nostrum aut dolorem internos et voluptatem aperiam sed deserunt iste. Eos voluptas libero et iusto velit eum perferendis sequi eos autem possimus qui iure numquam ut autem officiis.

        Aut tenetur veniam sed veritatis voluptatibus quo magnam similique hic voluptates natus? Aut explicabo alias et inventore veritatis ab quia assumenda sit excepturi maxime in assumenda labore qui dicta placeat eos rerum corrupti. Ut illo enim et assumenda sunt qui repellat corporis! Qui nulla illum et quisquam temporibus est laudantium itaque.
        
        Eos placeat consequatur ut consequatur inventore sed repudiandae ipsam non dolorem quas ex maiores magni non doloremque corporis. Eos commodi voluptatem vel iure consectetur ut voluptatibus eveniet id incidunt consequatur ad dicta dicta ea odio aliquam!`
      } </p>
    )
  }
  return (

    <ParallaxBanner
      layers={[image, text, paragraph]}
      className="parallax"
    />
  )
}

export default App;
