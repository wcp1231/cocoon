import { createSlice, current } from "@reduxjs/toolkit";
import {RootState} from "../../store";

export interface RecordMessage {
  id: number,
  captureTime: number,
  isRequest: boolean,
  meta: { [key: string]: string },
  header: { [key: string]: string },
  payload: { [key: string]: { [key: string]: any } },
  body: string,
}

export interface RecordItem {
  id: number,
  protocol: string,
  timespan: number,
  request: RecordMessage,
  response: RecordMessage | undefined,
}

export interface RecordsState {
  records: RecordItem[]
}

const initialState: RecordsState = {
  records: [],
};

const getProtocol = (request: RecordMessage): string => {
  let protocol = "UNKNOWN";
  if (!request) {
    return protocol;
  }
  return request.meta["PROTOCOL"] || protocol;
}

const appendRequest = (request: RecordMessage, records: RecordItem[]): RecordItem[] => {
  return [ ...records, {
    id: request.id,
    protocol: getProtocol(request),
    timespan: -1,
    request: request,
    response: undefined,
  }]
}

const updateResponse = (response: RecordMessage, record: RecordItem): RecordItem => {
  return {
    ...record,
    timespan: response.captureTime - record.request.captureTime,
    response: response,
  }
}

const appendResponse = (response: RecordMessage, records: RecordItem[]): RecordItem[] => {
  return records.map(record => {
    if (record.id === response.id) {
      return updateResponse(response, record);
    }
    return { ...record };
  })
}

export const recordsSlice = createSlice({
  name: 'records',
  initialState,
  reducers: {
    addRequest: (state, { payload }) => {
      state.records = appendRequest(payload, state.records);
    },
    addResponse: (state, { payload }) => {
      state.records = appendResponse(payload, state.records);
    },
  }
});

export const { addRequest, addResponse } = recordsSlice.actions;

export const selectRecords = (state: RootState): RecordItem[] => state.records.records;

export default recordsSlice.reducer;