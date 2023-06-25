import DevicesOtherIcon from '@mui/icons-material/DevicesOther';
import InsertEmoticonIcon from '@mui/icons-material/InsertEmoticon';
import AdminPanelSettingsIcon from '@mui/icons-material/AdminPanelSettings';
import CampaignIcon from '@mui/icons-material/Campaign';

const DeviceListIcon = ({ fontSize = 16 }) => {
  return (
    <>
      <DevicesOtherIcon style={{ fontSize: fontSize }} />
    </>
  );
};

const UserListIcon = ({ fontSize = 16 }) => {
  return (
    <>
      <InsertEmoticonIcon style={{ fontSize: fontSize }} />
    </>
  );
};

const AdminListIcon = ({ fontSize = 16 }) => {
  return (
    <>
      <AdminPanelSettingsIcon style={{ fontSize: fontSize }} />
    </>
  );
};


const Campaign = ({ fontSize = 16 }) => {
  return (
    <>
      <CampaignIcon style={{ fontSize: fontSize }} />
    </>
  );
};


export { DeviceListIcon, UserListIcon, AdminListIcon, Campaign };