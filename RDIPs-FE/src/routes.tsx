import { useRoutes, RouteObject } from 'react-router-dom';
import HomePage from './components/pages/admin/home/Home.page';
import TempPage from './components/pages/Temp.page';
import ListDevices from './components/pages/admin/list-devices/ListDevices';
import DetailDevice from './components/pages/admin/list-devices/detail-device/DetailDevice';
import DetailHistoryLog from './components/pages/admin/list-devices/detail-device/DetailHistoryLog';
import ListUsers from './components/pages/admin/list-users/ListUsers';
import DetailUser from './components/pages/admin/list-users/detail-user/DetailUser';
import Campaign from './components/pages/admin/campaign/Campaign';
import ListAdmin from './components/pages/admin/list-admin/ListAdmin';

export function AdminRoute(): ReturnType<typeof useRoutes> {
  const routes: RouteObject[] = [
    {
      path: '/',
      element: <TempPage />,
    },
    {
      path: '/admin',
      element: <HomePage />,
    },
    {
      path: '/list-devices',
      element: <ListDevices />,
    },
    {
      path: '/list-users',
      element: <ListUsers />,
    },
    {
      path: '/list-admin',
      element: <ListAdmin />,
    },
    {
      path: '/campaign',
      element: <Campaign />,
    },
    {
      path: '/detail-user',
      element: <DetailUser />,
    },
    {
      path: 'detail-history-log',
      element: <DetailHistoryLog />,
    },
    {
      path: '/detail-device',
      element: <DetailDevice />,
    },
  ];

  return useRoutes(routes);
}
