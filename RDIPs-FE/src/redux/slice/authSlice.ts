import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import authService from '../../service/auth.service';
import { AccountInfo } from '../../model/page';
import { GetThunkAPI } from '@reduxjs/toolkit/dist/createAsyncThunk';
const user = {};

export const login = createAsyncThunk(
  "auth/login",
  async (params : AccountInfo, thunkAPI: GetThunkAPI<any>) => {
    const data = await authService.login(params, thunkAPI);
    return data
  }
);

const initialState = user
  ? { isLoggedIn: true, user }
  : { isLoggedIn: false, user: null };

const authSlice = createSlice({
  name: "auth",
  initialState,
  reducers: {},
});

const { reducer } = authSlice;
export default reducer;
