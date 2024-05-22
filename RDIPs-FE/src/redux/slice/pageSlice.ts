import { createSlice } from '@reduxjs/toolkit';
import { Page } from '../../model/page';
const initialState = {};
const pageSlice = createSlice({
  name: 'page',
  initialState,
  reducers: {
    setNavbarTitle: (state, action) => {
      return { navbarTitle: action.payload };
    },
  },
});

export const getNavbarTitle = (state: { page: Page }) => {
  return state.page.navbarTitle;
};
const { reducer, actions } = pageSlice;
export const { setNavbarTitle } = actions;
export default reducer;
