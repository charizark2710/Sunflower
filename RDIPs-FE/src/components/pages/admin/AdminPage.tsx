import FooterOrganism from '../../organisms/common/footer/Footer.organism';
import HeaderOrganism from '../../organisms/common/header/Header.organism';
import HeaderTemplate from '../../templates/common/Header.template';
import './AdminPage.scss';

interface AdminPageProps {
  children: React.ReactNode;
}

function AdminPage(props: AdminPageProps) {
  return (
    <>
      <HeaderTemplate header={<HeaderOrganism />} />
      <div className='container'>{props.children}</div>
      <FooterOrganism />
    </>
  );
}

export default AdminPage;
