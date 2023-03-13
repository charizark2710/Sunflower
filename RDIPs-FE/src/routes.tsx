import { useRoutes, RouteObject } from "react-router-dom"
import ParallaxPage from "./components/pages/Parallax.page"
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

export function ScrollParallaxRoute(): ReturnType<typeof useRoutes> {
  const routes: RouteObject[] = [
    {
      path: "/parallax",
      element: <ParallaxPage />
    }
  ]

  return useRoutes(routes)
}