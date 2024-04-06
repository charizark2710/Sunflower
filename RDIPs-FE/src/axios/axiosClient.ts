import axios from 'axios';

const baseURL = process.env.REACT_APP_API_URL || "/api"

export const axiosClient = axios.create({
  baseURL,
  headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' + localStorage.getItem('token')
    },
})

export function setupAxios(store: any) {  
    axiosClient.interceptors.request.use(
      function(config: any) {
        return config
      },
      function(error: any) {
        return Promise.reject(error)
      },
    )

    axiosClient.interceptors.response.use(
      function(response: any) {
        return response.data
      },
      async function(error: any) {
        return Promise.reject(error)
      },
    )
  }
