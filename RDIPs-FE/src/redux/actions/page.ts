import { SET_TITLE } from '../constants/constant'
import { GET_TITLE} from '../constants/constant'
export function getPage() {
  return {
    type: GET_TITLE,
  }
}

export function setPage (title : string) {
  return {
    type: SET_TITLE,
    navbarTitle: title
  }
}