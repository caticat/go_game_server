package pnet

var (
	g_socketMgr PSocketManager = nil
	g_sessionID int64          = 0
)

func GetSocketMgr() PSocketManager          { return g_socketMgr }
func setSocketMgr(socketMgr PSocketManager) { g_socketMgr = socketMgr }
func GenSessionID() int64                   { g_sessionID += 1; return g_sessionID }
