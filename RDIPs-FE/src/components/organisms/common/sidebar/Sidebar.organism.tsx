import SunflowerLabel from '../../../molecules/label/SunflowerLabel.mocules';
import './Sidebar.organism.scss';

const sideBarItems = [
  <span>
    <strong>SUN</strong>flower
  </span>,
  'List Devices',
  'List Users',
  'Button name',
  'Button name',
  'Button name',
];

const topLabelStyle = {
  fontSize: '16px',
  lineHeight: '19px',
};

function SidebarOrganism({onClick, size} : {onClick: (args: any) => void, size: string}) {
  return (
    <div className='sidebar'>
      {
        sideBarItems.map((item, i) => {
          return i === 0 ? (
            <SunflowerLabel key={i} labelName={item} iconPos={1} height='62px' style={topLabelStyle} icon='toggle' onClick={onClick} size={size}/>
          ) : (
            <SunflowerLabel key={i} labelName={item} size={size} />
          );
        })
      }
      
    </div>
  );
}

export default SidebarOrganism;