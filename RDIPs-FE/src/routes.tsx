import { useRoutes, RouteObject } from "react-router-dom"
import HomePage from './components/pages/admin/home/Home.page'
import TempPage from './components/pages/Temp.page'

export function AdminRoute(): ReturnType<typeof useRoutes> {
  const routes: RouteObject[] = [
    {
      path: "/",
      element: <TempPage />
    },
    {
      path: "/admin",
      element: <HomePage />
    }
  ]

  return useRoutes(routes)
}