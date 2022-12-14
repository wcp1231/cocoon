import styled from "@emotion/styled";
import {RecordItem, RecordMessage} from "./recordsSlice";
import {Chip} from "@mui/material";

const StatusTag = styled.span`
  margin-right:1rem;
  font-size: 0.75rem;
  font-weight: 500;
  padding: 0.25rem 0.4rem;
  border-radius: 4px;
  
  color: ${props => props.color || '#000'};
  border: solid 1px ${props => props.color || '#000'};
`;

const HttpStatusTag = (response: RecordMessage, isMock: boolean) => {
  const statusCode = response.payload['HTTP_RESPONSE']['StatusCode'] as number;
  let color = '#28a745';
  if (statusCode >= 300) color = '#6c757d';
  if (statusCode >= 400) color = '#ffc107';
  if (statusCode >= 500) color = '#dc3545';
  return (
    <>
      <StatusTag color={color}>{statusCode}</StatusTag>
      {isMock && (<Chip variant="outlined" color="info" size="small" label="MOCK" />)}
    </>
  )
}

export default function(record: RecordItem) {
  if (record.timespan < 0) {
    return (<StatusTag>Pending</StatusTag>);
  }
  const response = record.response as RecordMessage;
  const isMock = response.meta['MOCK'] === 'true';
  if (record.protocol === 'HTTP') {
    return HttpStatusTag(response, isMock);
  }
  let color = "#28a745";
  let text = "OK";
  if (response.meta['STATUS'] !== 'ok') {
    color = "#dc3545";
    text = "Exception";
  }
  return (
    <>
      <StatusTag color={color}>{text}</StatusTag>
      {isMock && (<Chip variant="outlined" color="info" size="small" label="MOCK" />)}
    </>
  );
}