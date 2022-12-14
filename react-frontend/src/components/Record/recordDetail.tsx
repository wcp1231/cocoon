import {Card, CardContent, CardHeader, IconButton} from "@mui/material";
//import ClearIcon from '@mui/icons-material/Clear';
import React from "react";
import {RecordItem} from "./recordsSlice";
import DubboRecordDetail from "./dubboRecordDetail";
import HttpRecordDetail from "./httpRecordDetail";
import MongoRecordDetail from "./MongoRecordDetail";
import MysqlRecordDetail from "./MysqlRecordDetail";
import RedisRecordDetail from "./RedisRecordDetail";

const renderByProtocol = (record: RecordItem) => {
  const protocol = record.protocol;
  if (protocol === 'Dubbo') {
    return (<DubboRecordDetail record={record} />)
  }
  if (protocol === 'HTTP') {
    return (<HttpRecordDetail record={record} />)
  }
  if (protocol === 'Mongo') {
    return (<MongoRecordDetail record={record} />)
  }
  if (protocol === 'Mysql') {
    return (<MysqlRecordDetail record={record} />)
  }
  if (protocol === 'Redis') {
    return (<RedisRecordDetail record={record} />)
  }
  return (<div>{record.toString()}</div>)
}

export default function ({ record, onClose }: { record: RecordItem | undefined, onClose: () => void }) {
  if (!record) {
    return (<div></div>)
  }
  return (
    <Card>
      <CardHeader
        action={<IconButton onClick={onClose}>X</IconButton>}
        title="Detail"
      />
      <CardContent>
        {renderByProtocol(record)}
      </CardContent>
    </Card>
  )
}