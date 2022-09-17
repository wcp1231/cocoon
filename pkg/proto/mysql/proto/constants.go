package proto

/***************************************************/
// https://dev.mysql.com/doc/internals/en/command-phase.html
// include/my_command.h
const (
	COM_SLEEP               = 0x00
	COM_QUIT                = 0x01
	COM_INIT_DB             = 0x02
	COM_QUERY               = 0x03
	COM_FIELD_LIST          = 0x04
	COM_CREATE_DB           = 0x05
	COM_DROP_DB             = 0x06
	COM_REFRESH             = 0x07
	COM_SHUTDOWN            = 0x08
	COM_STATISTICS          = 0x09
	COM_PROCESS_INFO        = 0x0a
	COM_CONNECT             = 0x0b
	COM_PROCESS_KILL        = 0x0c
	COM_DEBUG               = 0x0d
	COM_PING                = 0x0e
	COM_TIME                = 0x0f
	COM_DELAYED_INSERT      = 0x10
	COM_CHANGE_USER         = 0x11
	COM_BINLOG_DUMP         = 0x12
	COM_TABLE_DUMP          = 0x13
	COM_CONNECT_OUT         = 0x14
	COM_REGISTER_SLAVE      = 0x15
	COM_STMT_PREPARE        = 0x16
	COM_STMT_EXECUTE        = 0x17
	COM_STMT_SEND_LONG_DATA = 0x18
	COM_STMT_CLOSE          = 0x19
	COM_STMT_RESET          = 0x1a
	COM_SET_OPTION          = 0x1b
	COM_STMT_FETCH          = 0x1c
	COM_DAEMON              = 0x1d
	COM_BINLOG_DUMP_GTID    = 0x1e
	COM_RESET_CONNECTION    = 0x1f
)

// CommandString used for translate cmd to string.
func CommandString(cmd byte) string {
	switch cmd {
	case COM_SLEEP:
		return "COM_SLEEP"
	case COM_QUIT:
		return "COM_QUIT"
	case COM_INIT_DB:
		return "COM_INIT_DB"
	case COM_QUERY:
		return "COM_QUERY"
	case COM_FIELD_LIST:
		return "COM_FIELD_LIST"
	case COM_CREATE_DB:
		return "COM_CREATE_DB"
	case COM_DROP_DB:
		return "COM_DROP_DB"
	case COM_REFRESH:
		return "COM_REFRESH"
	case COM_SHUTDOWN:
		return "COM_SHUTDOWN"
	case COM_STATISTICS:
		return "COM_STATISTICS"
	case COM_PROCESS_INFO:
		return "COM_PROCESS_INFO"
	case COM_CONNECT:
		return "COM_CONNECT"
	case COM_PROCESS_KILL:
		return "COM_PROCESS_KILL"
	case COM_DEBUG:
		return "COM_DEBUG"
	case COM_PING:
		return "COM_PING"
	case COM_TIME:
		return "COM_TIME"
	case COM_DELAYED_INSERT:
		return "COM_DELAYED_INSERT"
	case COM_CHANGE_USER:
		return "COM_CHANGE_USER"
	case COM_BINLOG_DUMP:
		return "COM_BINLOG_DUMP"
	case COM_TABLE_DUMP:
		return "COM_TABLE_DUMP"
	case COM_CONNECT_OUT:
		return "COM_CONNECT_OUT"
	case COM_REGISTER_SLAVE:
		return "COM_REGISTER_SLAVE"
	case COM_STMT_PREPARE:
		return "COM_STMT_PREPARE"
	case COM_STMT_EXECUTE:
		return "COM_STMT_EXECUTE"
	case COM_STMT_SEND_LONG_DATA:
		return "COM_STMT_SEND_LONG_DATA"
	case COM_STMT_CLOSE:
		return "COM_STMT_CLOSE"
	case COM_STMT_RESET:
		return "COM_STMT_RESET"
	case COM_SET_OPTION:
		return "COM_SET_OPTION"
	case COM_STMT_FETCH:
		return "COM_STMT_FETCH"
	case COM_DAEMON:
		return "COM_DAEMON"
	case COM_BINLOG_DUMP_GTID:
		return "COM_BINLOG_DUMP_GTID"
	case COM_RESET_CONNECTION:
		return "COM_RESET_CONNECTION"
	}
	return "UNKNOWN"
}

// https://dev.mysql.com/doc/internals/en/capability-flags.html
// include/mysql_com.h
const (
	// new more secure password
	CLIENT_LONG_PASSWORD = 1

	// Found instead of affected rows
	CLIENT_FOUND_ROWS = uint32(1 << 1)

	// Get all column flags
	CLIENT_LONG_FLAG = uint32(1 << 2)

	// One can specify db on connect
	CLIENT_CONNECT_WITH_DB = uint32(1 << 3)

	// Don't allow database.table.column
	CLIENT_NO_SCHEMA = uint32(1 << 4)

	// Can use compression protocol
	CLIENT_COMPRESS = uint32(1 << 5)

	// Odbc client
	CLIENT_ODBC = uint32(1 << 6)

	// Can use LOAD DATA LOCAL
	CLIENT_LOCAL_FILES = uint32(1 << 7)

	// Ignore spaces before '('
	CLIENT_IGNORE_SPACE = uint32(1 << 8)

	// New 4.1 protocol
	CLIENT_PROTOCOL_41 = uint32(1 << 9)

	// This is an interactive client
	CLIENT_INTERACTIVE = uint32(1 << 10)

	// Switch to SSL after handshake
	CLIENT_SSL = uint32(1 << 11)

	// IGNORE sigpipes
	CLIENT_IGNORE_SIGPIPE = uint32(1 << 12)

	// Client knows about transactions
	CLIENT_TRANSACTIONS = uint32(1 << 13)

	// Old flag for 4.1 protocol
	CLIENT_RESERVED = uint32(1 << 14)

	// Old flag for 4.1 authentication
	CLIENT_SECURE_CONNECTION = uint32(1 << 15)

	// Enable/disable multi-stmt support
	CLIENT_MULTI_STATEMENTS = uint32(1 << 16)

	// Enable/disable multi-results
	CLIENT_MULTI_RESULTS = uint32(1 << 17)

	// Multi-results in PS-protocol
	CLIENT_PS_MULTI_RESULTS = uint32(1 << 18)

	// Client supports plugin authentication
	CLIENT_PLUGIN_AUTH = uint32(1 << 19)

	// Client supports connection attributes
	CLIENT_CONNECT_ATTRS = uint32(1 << 20)

	//  Enable authentication response packet to be larger than 255 bytes
	CLIENT_PLUGIN_AUTH_LENENC_CLIENT_DATA = uint32(1 << 21)

	// Don't close the connection for a connection with expired password
	CLIENT_CAN_HANDLE_EXPIRED_PASSWORDS = uint32(1 << 22)

	// Capable of handling server state change information. Its a hint to the
	// server to include the state change information in Ok packet.
	CLIENT_SESSION_TRACK = uint32(1 << 23)

	//Client no longer needs EOF packet
	CLIENT_DEPRECATE_EOF = uint32(1 << 24)
)
