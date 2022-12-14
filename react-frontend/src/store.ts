import {configureStore} from "@reduxjs/toolkit";
import recordsReducer from "./components/Record/recordsSlice";

export const store = configureStore({
  reducer: {
    records: recordsReducer,
  }
});

export type AppDispatch = typeof store.dispatch;
export type RootState = ReturnType<typeof store.getState>;