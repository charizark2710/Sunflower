import { FETCH_POST_ERROR, FETCH_POST_REQUEST, FETCH_POST_SUCCESS } from '../constants/post';


export const loadPosts = () => async (dispatch : any) => {
  try {
      dispatch({ type: FETCH_POST_REQUEST });

      const url = "https://jsonplaceholder.typicode.com/posts";
      const response = await fetch(url)
      const responseBody = await response.json();
      dispatch({
          type: FETCH_POST_SUCCESS,
          data: responseBody
      });
  } catch (error) {
      console.error(error);
      dispatch({
          type: FETCH_POST_ERROR,
          message: error
      });
  }
}