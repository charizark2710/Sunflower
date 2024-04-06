import { configureStore } from '@reduxjs/toolkit';
import {
  useDispatch as useReduxDispatch,
  useSelector as useReduxSelector,
} from "react-redux";
import pageReducers from './slice/pageSlice';
import authReducers from './slice/authSlice';

const reducer = {
  page: pageReducers,
  auth: authReducers
};

export const useSelector = useReduxSelector;
export const useDispatch = () => useReduxDispatch();

const store = configureStore({
  reducer: reducer,
  devTools: true,
});

export default store;