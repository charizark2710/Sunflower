import SunflowerLabel from '../../../molecules/label/SunflowerLabel.mocules';
import './Sidebar.organism.scss';

const sideBarItems = [
  {
    to: '/',
    className: 'link-item',
    children: (
      <span>
        <strong>SUN</strong>flower
      </span>
    ),
  },
  { to: '/list-devices', className: 'link-item', children: 'List Devices' },
  { to: '/', className: 'link-item', children: 'List Users' },
  { to: '/', className: 'link-item', children: 'Button name' },
  { to: '/', className: 'link-item', children: 'Button name' },
  { to: '/', className: 'link-item', children: 'Button name' },
];

const topLabelStyle = {
  fontSize: '16px',
  lineHeight: '19px',
};

function SidebarOrganism({ onClick, size }: { onClick: (args: any) => void; size: string }) {
  return (
    <div className='sidebar'>
      {sideBarItems.map((item, i) => {
        return i === 0 ? (
          <SunflowerLabel
            key={i}
            link={item}
            iconPos={1}
            height='62px'
            style={topLabelStyle}
            icon='toggle'
            onClick={onClick}
            size={size}
          />
        ) : (
          <SunflowerLabel key={i} link={item} size={size} />
        );
      })}
    </div>
  );
}

export default SidebarOrganism;