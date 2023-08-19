package reply

type UnKnownErrReply struct {
}

var unknownErrBytes = []byte("-Err unknown\r\n")

func (u *UnKnownErrReply) Error() string {
	return "Err unknown"
}

func (u *UnKnownErrReply) ToBytes() []byte {
	return unknownErrBytes
}

type ArgNumErrReply struct {
	Cmd string
}

func (a *ArgNumErrReply) Error() string {
	return "-ERR wrong number of arguments for '" + a.Cmd + "' command"
}

func (a *ArgNumErrReply) ToBytes() []byte {
	return []byte("-ERR wrong number of arguments for '" + a.Cmd + "' command\r\n")
}

func NewArgNumErrReply(cmd string) *ArgNumErrReply {
	return &ArgNumErrReply{Cmd: cmd}
}

type SyntaxErrReply struct{}

var syntaxErrBytes = []byte("-Err syntax error\r\n")

func (s *SyntaxErrReply) Error() string {
	return "Err syntax error"
}

func (s *SyntaxErrReply) ToBytes() []byte {
	return syntaxErrBytes
}

type WrongTypeErrReply struct{}

var wrongTypeErrBytes = []byte("-Err WRONGTYPE Operation against a holding the wrong kind of value\r\n")

func (w *WrongTypeErrReply) Error() string {
	return "Err WRONGTYPE Operation against a holding the wrong kind of value"
}

func (w *WrongTypeErrReply) ToBytes() []byte {
	return wrongTypeErrBytes
}

type ProtocolErrReply struct {
	Msg string
}

func (p *ProtocolErrReply) Error() string {
	return "Err Protocol err: '" + p.Msg
}

func (p *ProtocolErrReply) ToBytes() []byte {
	return []byte("-Err Protocol err: '" + p.Msg + "'\r\n")
}
