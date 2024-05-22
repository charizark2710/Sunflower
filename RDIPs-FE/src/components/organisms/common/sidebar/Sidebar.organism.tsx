import { Box, Typography } from '@mui/material'
import {
  AdminListIcon,
  Campaign,
  DeviceListIcon,
  UserListIcon,
} from "../../../atoms/icon/ListIcon.atom"
import SunFlowerLogo from "../../../atoms/image/SunFlowerLogo.atom"
import { LogoutButtonMolecules } from "../../../molecules/clutter-button-mocules/LogoutButton.mocules"
import SunflowerLabel from "../../../molecules/label/SunflowerLabel.mocules"
import "./Sidebar.organism.scss"

const topLabelStyle = {
  fontSize: "16px",
  lineHeight: "19px",
  marginBottom: "20px",
}

function SidebarOrganism({
  onClick,
  changeStateSideBar,
  size,
  collapse,
}: {
  onClick?: (args: any) => void
  changeStateSideBar?: (args: any) => void
  size: string
  collapse: boolean
}) {
  const sideBarItems = [
    {
      to: "/",
      className: "link-item",
      children: (
        <Typography component={'span'} style={{ fontSize: "16px" }}>
          <SunFlowerLogo w={"100px"} label={!collapse} />
        </Typography>
      ),
    },
    {
      to: "/list-devices",
      className: "link-item",
      children: "Device's List",
      icon: <DeviceListIcon />,
    },
    {
      to: "/list-users",
      className: "link-item",
      children: "User's List",
      icon: <UserListIcon />,
    },
    {
      to: "/list-admin",
      className: "link-item",
      children: "Admin's List",
      icon: <AdminListIcon />,
    },
    {
      to: "/campaign",
      className: "link-item",
      children: "Campaign's List",
      icon: <Campaign />,
    },
  ]

  return (
    <Box className="sidebar">
      <Box>
        {sideBarItems.map((item, i) => {
          return i === 0 ? (
            <Box key={i}>
              <SunflowerLabel
                key={item.to + i}
                link={item}
                iconPos={1}
                height="32px"
                style={topLabelStyle}
                icon="toggle"
                onClick={onClick}
                size={size}
                isHomepage={true}
              />
              <Box className="search-area" onClick={changeStateSideBar}></Box>
            </Box>
          ) : (
            <SunflowerLabel
              state={collapse}
              key={item.to + i}
              link={item}
              size={size}
              children={""}
              specialIcon={item.icon}
            />
          )
        })}
      </Box>
      <LogoutButtonMolecules />
    </Box>
  )
}

export default SidebarOrganism
