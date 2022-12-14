import React from "react";
import {
  ColumnDef, flexRender, getCoreRowModel, Row, useReactTable
} from "@tanstack/react-table";
import { useVirtual } from "react-virtual";
import {useAppDispatch, useAppSelector} from "../../hooks";
import {addRequest, selectRecords, RecordItem, addResponse, RecordMessage} from "./recordsSlice";
import useWebSocket from "react-use-websocket";
import timespan from "./timespan";
import briefRequest from "./briefRequest";
import briefResponse from "./briefResponse";
import {
  Container,
  Drawer,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow
} from "@mui/material";
import ProtocolTag from "./protocolTag";
import RecordDetail from "./recordDetail";

//const wsUrl = `ws://${window.location.host}/api/ws`;
const wsUrl = `ws://127.0.0.1:7070/api/ws`;

export default function RecordView() {

  const dispatch = useAppDispatch();

  // TODO 打开页面就建立 WS 监听
  const { sendJsonMessage, getWebSocket } = useWebSocket(wsUrl, {
    onOpen: () => console.log('WebSocket connection opened'),
    onClose: () => console.log('WebSocket connection closed'),
    shouldReconnect: (closeEvent) => true,
    onMessage: (event: WebSocketEventMap["message"]) => processMessage(event)
  });

  const processMessage = (event: { data: string }) => {
    const record = JSON.parse(event.data)
    if (record.isRequest) {
      dispatch(addRequest(record))
    } else {
      dispatch(addResponse(record))
    }
  }

  const columns = React.useMemo<ColumnDef<RecordItem>[]>(
    () => [
      { accessorKey: 'id', header: 'ID', size: 60, },
      {
        accessorKey: 'protocol', header: 'Protocol', size: 120,
        cell: info => ProtocolTag(info.getValue() as string),
      },
      {
        accessorKey: 'request',
        header: 'Request',
        cell: info => briefRequest(info.row.getValue('protocol') as string, info.getValue() as RecordMessage),
      },
      {
        accessorKey: 'response',
        header: 'Response',
        cell: info => briefResponse(info.row.original as RecordItem),
      },
      {
        accessorKey: 'timespan',
        header: 'Time',
        cell: info => timespan(info.getValue() as number),
        size: 120,
      }
    ],
    []
  );

  const records: RecordItem[] = useAppSelector(selectRecords)
  const [drawerState, setDrawerState] = React.useState(false)
  const [selectedRecord, setSelectedRecord] = React.useState<RecordItem>()
  const toggleDrawer = (open: boolean) => setDrawerState(open)
  const selectRecord = (record: RecordItem) => {
    setSelectedRecord(record);
    toggleDrawer(true);
  }

  return (
    <>
      <LocalTable data={records} columns={columns} onRowClick={selectRecord}/>
      <Drawer
        anchor="right"
        variant="persistent"
        hideBackdrop={true}
        open={drawerState}
      >
        <RecordDetail record={selectedRecord} onClose={() => toggleDrawer(false)} />
      </Drawer>
    </>
  )
}

function LocalTable(
  { data, columns, onRowClick }: { data: RecordItem[], columns: ColumnDef<RecordItem>[], onRowClick: (r: RecordItem) => void }
) {
  const table = useReactTable({
    defaultColumn: { minSize: 0, size: 0 },
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
    debugTable: true,
  });
  const tableContainerRef = React.useRef<HTMLDivElement>(null);

  const { rows } = table.getRowModel()
  const rowVirtualizer = useVirtual({
    parentRef: tableContainerRef,
    size: rows.length,
    overscan: 10,
  });
  const { virtualItems: virtualRows, totalSize } = rowVirtualizer;

  const paddingTop = virtualRows.length > 0 ? virtualRows?.[0]?.start || 0 : 0;
  const paddingBottom = virtualRows.length > 0 ? totalSize - (virtualRows?.[virtualRows.length-1]?.end || 0) : 0;

  return (
    <TableContainer component={Paper}>
      <Table sx={{ minWidth: 650 }}>
        <TableHead>
        {table.getHeaderGroups().map(headerGroup => (
          <TableRow key={headerGroup.id}>
          {headerGroup.headers.map(header => {
            return (
              <TableCell key={header.id} style={{ width: header.getSize() || undefined }}>
                {flexRender(header.column.columnDef.header, header.getContext())}
              </TableCell>
            )
          })}
          </TableRow>
        ))}
        </TableHead>
        <TableBody>
        {paddingTop > 0 && (
          <TableRow>
            <TableCell style={{ height: `${paddingTop}` }} />
          </TableRow>
        )}
        {table.getRowModel().rows.map(row => {
          return (
            <TableRow key={row.id} onClick={() => { onRowClick(row.original) }}>
            {row.getVisibleCells().map(cell => {
              return (
                <TableCell key={cell.id}>
                  {flexRender(cell.column.columnDef.cell, cell.getContext())}
                </TableCell>
              )
            })}
            </TableRow>
          )
        })}
        {paddingBottom > 0 && (
          <TableRow>
            <TableCell style={{ height: `${paddingBottom}` }} />
          </TableRow>
        )}
        </TableBody>
      </Table>
    </TableContainer>
  )
}

function RawTable({ data, columns }: { data: RecordItem[], columns: ColumnDef<RecordItem>[] }) {
  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
    debugTable: true,
  });

  const tableContainerRef = React.useRef<HTMLDivElement>(null);

  const { rows } = table.getRowModel()
  const rowVirtualizer = useVirtual({
    parentRef: tableContainerRef,
    size: rows.length,
    overscan: 10,
  });
  const { virtualItems: virtualRows, totalSize } = rowVirtualizer;

  const paddingTop = virtualRows.length > 0 ? virtualRows?.[0]?.start || 0 : 0;
  const paddingBottom = virtualRows.length > 0 ? totalSize - (virtualRows?.[virtualRows.length-1]?.end || 0) : 0;

  return (
    <Container maxWidth={false}>
      <div>Record View</div>
      <div ref={tableContainerRef} className="table-container">
        <table className="w-full">
          <thead>
          {table.getHeaderGroups().map(headerGroup => (
            <tr key={headerGroup.id}>
              {headerGroup.headers.map(header => {
                return (
                  <th key={header.id} colSpan={header.colSpan} style={{ width: header.getSize() }}>
                    {header.isPlaceholder ? null : (
                      <div>
                        {flexRender(header.column.columnDef.header, header.getContext())}
                      </div>
                    )}
                  </th>
                )
              })}
            </tr>
          ))}
          </thead>
          <tbody>
          {paddingTop > 0 && (
            <tr>
              <td style={{ height: `${paddingTop}` }} />
            </tr>
          )}
          {virtualRows.map(virtualRow => {
            const row = rows[virtualRow.index] as Row<RecordItem>;
            return (
              <tr key={row.id}>
                {row.getVisibleCells().map(cell => {
                  return (
                    <td key={cell.id}>
                      {flexRender(cell.column.columnDef.cell, cell.getContext())}
                    </td>
                  )
                })}
              </tr>
            )
          })}
          {paddingBottom > 0 && (
            <tr>
              <td style={{ height: `${paddingBottom}` }} />
            </tr>
          )}
          </tbody>
        </table>
      </div>
    </Container>
  )
}