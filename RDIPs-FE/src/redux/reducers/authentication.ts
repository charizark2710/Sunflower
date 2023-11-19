import { GET_AUTHEN, SET_AUTHEN } from '../constants/constant';

export default function authenticate(state = false, action = {} as any) {
  switch (action.type) {
    case GET_AUTHEN:
      return state;
    case SET_AUTHEN:
      return {
        isLogin: action.isLogin
      };
    default:
      return null;
  }
}
