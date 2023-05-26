const initialState = {
  requesting: false,
  success: false,
  message: null,
  data: null
}

function exampleReducers(state = initialState, payload : any) {
  switch (payload.type) {
      case "FETCH_EXAMPLE_REQUEST":
          return {
              ...state,
              requesting: true
          };
      case "FETCH_EXAMPLE_SUCCESS":
          return {
              ...state,
              requesting: false,
              success: true,
              data: payload.data
          };
      case "FETCH_EXAMPLE_ERROR":
          return {
              ...state,
              requesting: false,
              message: payload.message
          };

      default:
          return state;
  }
}

export default exampleReducers.bind(this);