import styled from "@emotion/styled";
import {RecordMessage} from "./recordsSlice";

const RequestBriefTag = styled.span`
  margin-right: 1rem;
  font-size: 0.75rem;
  padding: 0.25rem 0.4rem;
  border: solid 1px black;
  border-radius: 4px;
`;

export default function(protocol: string, request: RecordMessage) {
  if (protocol === 'Dubbo') {}
  if (protocol === 'HTTP') {
    return (
      <>
        <RequestBriefTag>{ request.payload['HTTP_REQUEST']['Method'].toString() }</RequestBriefTag>
        <span>{request.payload['HTTP_REQUEST']['Host'].toString()} {request.payload['HTTP_REQUEST']['URL'].toString()}</span>
      </>
    )
  }
  if (protocol === 'Mongo') {
  }
  if (protocol === 'Mysql') {
    return (
      <>
        <RequestBriefTag>{request.payload["MYSQL_OP_TYPE"].toString()}</RequestBriefTag>
        <span className="request-column-info">{request.payload["MYSQL_QUERY"].toString()}</span>
      </>
    )
  }
  if (protocol === 'Redis') {
  }
  return (<div>{protocol} {request.id}</div>)
}