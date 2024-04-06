import { GetThunkAPI } from '@reduxjs/toolkit/dist/createAsyncThunk';
import { axiosClient } from '../axios/axiosClient';
import { AccountInfo } from '../model/page';

const URL = '/login';

const login = async (params: AccountInfo, thunkAPI: GetThunkAPI<any>) => {
  try {
    const res = await axiosClient.post(URL, {username: 'admin', password: 'admin'});
    localStorage.setItem(
      'token',
      res.data.access_token ? res.data.access_token : 'eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJIaTd1cnpqT2phY3lnbEdjYzVkMWlUZ0tiZWVmY21oa2JXNkF0REtZMEw4In0.eyJleHAiOjE3MTIyNTEzOTIsImlhdCI6MTcxMjI1MTMzMiwianRpIjoiOGE4YWU5MTYtMDU5ZS00ZjYxLTgwNGUtNDhmZmFiZDMwOGYwIiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDo5MDkxL3JlYWxtcy9tYXN0ZXIiLCJzdWIiOiJhMWFmYmM1ZS0zMTk5LTRhZjctOGQyMi1iNzM2YTI4NWQ2MWMiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJhZG1pbi1jbGkiLCJzZXNzaW9uX3N0YXRlIjoiNDA3Y2JmZmEtNDZhYy00NjU3LTk2MzgtMzBkNDE5MDMzNzRlIiwiYWNyIjoiMSIsInNjb3BlIjoicHJvZmlsZSBlbWFpbCIsInNpZCI6IjQwN2NiZmZhLTQ2YWMtNDY1Ny05NjM4LTMwZDQxOTAzMzc0ZSIsImVtYWlsX3ZlcmlmaWVkIjpmYWxzZSwicHJlZmVycmVkX3VzZXJuYW1lIjoibHludHQifQ.LEia_2yuor3uMLfXPRKEfC7xls7ET8rCGZP6bqMsV6ncnh3KtpWJdUBNCvxfaeW5luGkeNjKEL5wRuzCku-BVoceErOsoeKY-A7jbxAyv8ZFImhKxtsOeCz7pwPqFPwotmwfVKZMrwpIT4KzTvKp6WMjqc-q2L6MhlrEtvx8Tk4unMbVgMiJP0oTQfciYRjmlQtdfzVga_Mi2KkHFWxzI9jPI_PocB3KEAALeBKEhTU8OWvMSy4mLj61JHzivyRrPaYKsG6xlyjypima824QuTaYqUkf8u1p98AcgeM7Z4y56ObFzMu0RnAHqAhQ_K9O1_z09K-39ohmklgY79LYnQ'
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
