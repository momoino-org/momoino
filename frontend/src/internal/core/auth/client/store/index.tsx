import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { U } from 'ts-toolbelt';
import { Profile } from '../../shared';

interface AuthState {
  profile: U.Nullable<Profile>;
}

const initialState: AuthState = {
  profile: null,
};

export const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    setProfile: (state, action: PayloadAction<U.Nullable<Profile>>) => {
      state.profile = action.payload;
    },
  },
});
