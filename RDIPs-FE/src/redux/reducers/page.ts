import { GET_TITLE, SET_TITLE } from '../constants/constant';

export default function navbarTitle(state = null, action = {} as any) {
  switch (action.type) {
    case GET_TITLE:
      return state;
    case SET_TITLE:
      return action.navbarTitle;
    default:
      return null;
  }
}
