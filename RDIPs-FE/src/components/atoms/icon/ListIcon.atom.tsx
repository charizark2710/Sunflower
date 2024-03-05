import sidebar_admin from '../../../assets/icons/sidebar_admin.svg';
import sidebar_campaign from '../../../assets/icons/sidebar_campaign.svg';
import sidebar_device from '../../../assets/icons/sidebar_device.svg';
import sidebar_user from '../../../assets/icons/sidebar_user.svg';
import sidebar_logout from '../../../assets/icons/sidebar_logout.svg';

const LogoutIcon = ({ fontSize = 16 }) => {
  return (
    <img src={sidebar_logout} style={{ fontSize: fontSize, verticalAlign: 'middle' }} />
  );
};

const DeviceListIcon = ({ fontSize = 16 }) => {
  return (
    <img src={sidebar_device} style={{ fontSize: fontSize, verticalAlign: 'middle' }} />
  );
};

const UserListIcon = ({ fontSize = 16 }) => {
  return (
    <img src={sidebar_user} style={{ fontSize: fontSize, verticalAlign: 'middle' }} />
  );
};

const AdminListIcon = ({ fontSize = 16 }) => {
  return (
    <img src={sidebar_admin} style={{ fontSize: fontSize, verticalAlign: 'middle' }} />
  );
};


const Campaign = ({ fontSize = 16 }) => {
  return (
    <img src={sidebar_campaign} style={{ fontSize: fontSize, verticalAlign: 'middle' }} />
  );
};

export { AdminListIcon, Campaign, DeviceListIcon, UserListIcon, LogoutIcon };