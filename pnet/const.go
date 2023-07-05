package pnet

type messageHandler_t func(*PRecvData) bool

const PSocket_ChanLen = 10 // 收发消息阻塞长度
