import React from "react";
import {RecordItem} from "./recordsSlice";

export default function ({ record }: { record: RecordItem | undefined }) {
  return (<div>HTTP {record?.request.header.toString()}</div>)
}