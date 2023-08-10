import { AdminListIcon, Campaign, DeviceListIcon, UserListIcon } from '../../../atoms/icon/ListIcon.atom';
import { StraightAtom } from '../../../atoms/straight/Straight.atom';
import SunflowerLabel from '../../../molecules/label/SunflowerLabel.mocules';
// import Image from '../../../atoms/image/Image';
// import sunflower from '../../../../assets/Sunflower.png';
import './Sidebar.organism.scss';

const sideBarItems = [
  {
    to: '/',
    className: 'link-item',
    children: (
      <span style={{ fontSize: '16px' }}>
        <strong>SUN</strong>flower
      </span>
    ),
  },
  { to: '/list-devices', className: 'link-item', children: 'List Devices', icon: <DeviceListIcon /> },
  { to: '/list-users', className: 'link-item', children: 'List Users', icon: <UserListIcon /> },
  { to: '/list-admin', className: 'link-item', children: 'List Admin', icon: <AdminListIcon /> },
  { to: '/campaign', className: 'link-item', children: 'Campaign', icon: <Campaign />},
];

const topLabelStyle = {
  fontSize: '16px',
  lineHeight: '19px',
  marginBottom: '20px',
};

const straight = <StraightAtom width='100%' thick='0.1px' color='#e4e4e4' />;

function SidebarOrganism({ onClick, size }: { onClick: (args: any) => void; size: string }) {
  return (
    <div className='sidebar'>
      {sideBarItems.map((item, i) => {
        return i === 0 ? (
          <div key={i}>
            <SunflowerLabel
              key={item.to + i}
              link={item}
              iconPos={1}
              height='10px'
              style={topLabelStyle}
              icon='toggle'
              onClick={onClick}
              size={size}
            />
            <div className = 'flex-justify-center' style={{marginBottom: '10px', height: '90px'}}> 

              {/* <Image url={sunflower} w ='50%' /> */}
            </div>
            {size === 'md' ? straight : ''}
          </div>
        ) : (
          <SunflowerLabel key={item.to + i} link={item} size={size} children={straight} specialIcon={item.icon} />
        );
      })}
    </div>
  );
}

export default SidebarOrganism;
