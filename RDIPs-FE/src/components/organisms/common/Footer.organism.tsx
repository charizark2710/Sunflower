import { ButtonAtom } from '../../atoms/button/Button.atom';
import { LinkMolecules } from '../../molecules/link-molecules/Link.molecules';
import './Footer.organism.scss';

function FooterOrganism() {
  return (
    <div className='footer-container'>
      <section className='footer-subscription'>
        <p className='footer-subscription-heading'>
          Join the Adventure newsletter to receive our best vacation deals
        </p>
        <p className='footer-subscription-text'>
          You can unsubscribe at any time.
        </p>
        <div className='input-areas'>
          <form>
            <input
              className='footer-input'
              name='email'
              type='email'
              placeholder='Your Email'
            />
            <ButtonAtom buttonStyle='btn--outline'>Subscribe</ButtonAtom>
          </form>
        </div>
      </section>
      <div className='footer-links'>
        <div className='footer-link-wrapper'>
          <div className='footer-link-items'>
            <h2>About Us</h2>
            <LinkMolecules to='/sign-up'>How it works</LinkMolecules>
            <LinkMolecules to='/'>Testimonials</LinkMolecules>
            <LinkMolecules to='/'>Careers</LinkMolecules>
            <LinkMolecules to='/'>Investors</LinkMolecules>
            <LinkMolecules to='/'>Terms of Service</LinkMolecules>
          </div>
          <div className='footer-link-items'>
            <h2>Contact Us</h2>
            <LinkMolecules to='/'>Contact</LinkMolecules>
            <LinkMolecules to='/'>Support</LinkMolecules>
            <LinkMolecules to='/'>Destinations</LinkMolecules>
            <LinkMolecules to='/'>Sponsorships</LinkMolecules>
          </div>
        </div>
        <div className='footer-link-wrapper'>
          <div className='footer-link-items'>
            <h2>Videos</h2>
            <LinkMolecules to='/'>Submit Video</LinkMolecules>
            <LinkMolecules to='/'>Ambassadors</LinkMolecules>
            <LinkMolecules to='/'>Agency</LinkMolecules>
            <LinkMolecules to='/'>Influencer</LinkMolecules>
          </div>
          <div className='footer-link-items'>
            <h2>Social Media</h2>
            <LinkMolecules to='/'>Instagram</LinkMolecules>
            <LinkMolecules to='/'>Facebook</LinkMolecules>
            <LinkMolecules to='/'>Youtube</LinkMolecules>
            <LinkMolecules to='/'>Twitter</LinkMolecules>
          </div>
        </div>
      </div>
      <section className='social-media'>
        <div className='social-media-wrap'>
          <div className='footer-logo'>
            <LinkMolecules to='/' className='social-logo'>
              RDIPS
              <i className='fab fa-typo3' />
            </LinkMolecules>
          </div>
          <small className='website-rights'>RDIPS © 2020</small>
          <div className='social-icons'>
            <LinkMolecules
              className='social-icon-link facebook'
              to='/'
              target='_blank'
              aria-label='Facebook'
            >
              <i className='fab fa-facebook-f' />
            </LinkMolecules>
            <LinkMolecules
              className='social-icon-link instagram'
              to='/'
              target='_blank'
              aria-label='Instagram'
            >
              <i className='fab fa-instagram' />
            </LinkMolecules>
            <LinkMolecules
              className='social-icon-link youtube'
              to='/'
              target='_blank'
              aria-label='Youtube'
            >
              <i className='fab fa-youtube' />
            </LinkMolecules>
            <LinkMolecules
              className='social-icon-link twitter'
              to='/'
              target='_blank'
              aria-label='Twitter'
            >
              <i className='fab fa-twitter' />
            </LinkMolecules>
            <LinkMolecules
              className='social-icon-link twitter'
              to='/'
              target='_blank'
              aria-label='LinkedIn'
            >
              <i className='fab fa-linkedin' />
            </LinkMolecules>
          </div>
        </div>
      </section>
    </div>
  );
}

export default FooterOrganism;