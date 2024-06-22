import { RouteObject, useRoutes } from 'react-router-dom';
import AdminPage from './components/pages/admin/AdminPage';
import Campaign from './components/pages/admin/campaign/Campaign';
import ListAdmin from './components/pages/admin/list-admin/ListAdmin';
import ListDevices from './components/pages/admin/list-devices/ListDevices';
import DetailDevice from './components/pages/admin/list-devices/detail-device/DetailDevice';
import DetailHistoryLog from './components/pages/admin/list-devices/detail-device/DetailHistoryLog';
import ListUsers from './components/pages/admin/list-users/ListUsers';
import DetailUser from './components/pages/admin/list-users/detail-user/DetailUser';
import LoginPage from './components/pages/authentication/LoginPage.page';
import Register from './components/pages/authentication/Register.page';
import LandingPage from './components/pages/landing/Landing.page';

export function CommonRoute(): ReturnType<typeof useRoutes> {
  const routes: RouteObject[] = [
    {
      path: '/',
      element: <AdminPage children={<ListDevices />} />
    },
    {
      path: '/login',
      element: <LoginPage />,
    },
    {
      path: '/register',
      element: <Register />,
    },
    {
      path: 'admin',
      element: <AdminPage />,
    },
    {
      path: '/list-devices',
      element: <AdminPage children={<ListDevices />} />
    },
    {
      path: '/list-users',
      element:  <AdminPage children={<ListUsers />} />,
    },
    {
      path: '/list-admin',
      element:  <AdminPage children={<ListAdmin />} />,
    },
    {
      path: '/campaign',
      element:  <AdminPage children={<Campaign />} />,
    },
    {
      path: '/detail-user',
      element: <AdminPage children={<DetailUser />} />,
    },
    {
      path: 'detail-history-log',
      element: <AdminPage children={<DetailHistoryLog />} />,
    },
    {
      path: '/detail-device',
      element: <AdminPage children={<DetailDevice />} />,
    },
  ];

  return useRoutes(routes);
}

export function AdminRoute(): ReturnType<typeof useRoutes> {
  const routes: RouteObject[] = [];

  return useRoutes(routes);
}
