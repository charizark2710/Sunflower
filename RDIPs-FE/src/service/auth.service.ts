import { GetThunkAPI } from '@reduxjs/toolkit/dist/createAsyncThunk';
import { axiosClientWithoutAuthentication } from '../axios/axiosClient';
import { AccountInfo } from '../model/page';

const URL = '/login';

const login = async (params: AccountInfo, thunkAPI: GetThunkAPI<any>) => {
  try {
    const res = await axiosClientWithoutAuthentication.post(URL, {username: params.username, password: params.password});
    localStorage.setItem(
      'token',
      res.data.data.access_token
    );
    return res;
  } catch (error) {
    return thunkAPI.rejectWithValue(0);
  }
};

const authService = {
  login,
};
export default authService;
