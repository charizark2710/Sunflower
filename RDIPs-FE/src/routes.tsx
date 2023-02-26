import { useRoutes, RouteObject } from "react-router-dom"
import TempPage from './components/pages/Temp.page'

export function AdminRoute(): ReturnType<typeof useRoutes> {
  const routes: RouteObject[] = [
    {
      path: "/",
      element: <TempPage />
    },
    {
      path: "/admin",
      element: <TempPage />
    }
  ]

  return useRoutes(routes)
}