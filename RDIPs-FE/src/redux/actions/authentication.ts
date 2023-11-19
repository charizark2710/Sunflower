import { GET_AUTHEN, SET_AUTHEN } from '../constants/constant'

export function getAuthen() {
  return {
    type: GET_AUTHEN,
  }
}

export function login(status : boolean) {
  return {
    type: SET_AUTHEN,
    isLogin: status
  }
}