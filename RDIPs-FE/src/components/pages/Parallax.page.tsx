import { ParallaxProvider } from 'react-scroll-parallax';
import BannerOneOrganism from '../organisms/parallax/BannerOne.organism';
import * as style from './Parallax.module.scss';
function ParallaxPage() {
    return (
        <ParallaxProvider>
          <BannerOneOrganism/>
        </ParallaxProvider>
    )
}

export default ParallaxPage
