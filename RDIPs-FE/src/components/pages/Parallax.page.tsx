import { ParallaxProvider } from 'react-scroll-parallax';
import BannerOneOrganism from '../organisms/parallax/BannerOne.organism';
import './Parallax.page.scss';
function ParallaxPage() {
    return (
      <div id='parallax'>
        <ParallaxProvider>
          <BannerOneOrganism/>
        </ParallaxProvider>
      </div>
    )
}

export default ParallaxPage
