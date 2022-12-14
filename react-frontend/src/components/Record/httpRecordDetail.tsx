import React from "react";
import {RecordItem} from "./recordsSlice";
import {Accordion, AccordionDetails, AccordionSummary, Divider, Typography} from "@mui/material";

const httpMethod = (record: RecordItem) => record.request.payload['HTTP_REQUEST']['Method'];
const httpStatus = (record: RecordItem) => record.response?.payload['HTTP_RESPONSE']['StatusCode'];
const httpURL = (record: RecordItem) => {
  const host = record.request.payload['HTTP_REQUEST']['Host'];
  const path = record.request.payload['HTTP_REQUEST']['URL'];
  return host + path;
}
const httpRequestHeaders = (record: RecordItem) => {
  return Object.entries(record.request.payload['HTTP_REQUEST']['Header']);
}
const httpResponseHeaders = (record: RecordItem) => {
  return Object.entries(record.response?.payload['HTTP_RESPONSE']['Header']);
}

export default function ({ record }: { record: RecordItem }) {
  return (
    <div>
      <Accordion>
        <AccordionSummary id="http-request">
          <Typography>Request</Typography>
        </AccordionSummary>
        <AccordionDetails>
          <p><strong>Method</strong>: {httpMethod(record)}</p>
          <p><strong>URL</strong>: {httpURL(record)}</p>
          <Divider textAlign="left">Request Header</Divider>
          {httpRequestHeaders(record).map((entry) => {
            return (
              <p><strong>{entry[0]}</strong>: {entry[1]?.toString()}</p>
            )
          })}
        </AccordionDetails>
      </Accordion>
      <Accordion>
        <AccordionSummary id="http-response">
          <Typography>Response</Typography>
        </AccordionSummary>
        <AccordionDetails>
          <p><strong>Status</strong>: {httpStatus(record)}</p>
          <Divider textAlign="left">Response Header</Divider>
          {httpResponseHeaders(record).map((entry) => {
            return (
              <p><strong>{entry[0]}</strong>: {entry[1]?.toString()}</p>
            )
          })}
        </AccordionDetails>
      </Accordion>
    </div>
  )
}