package parser

import (
	"bufio"
	"errors"
	"io"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/sunflower10086/go-redis/interface/resp"
	"github.com/sunflower10086/go-redis/lib/logger"
	"github.com/sunflower10086/go-redis/resp/reply"
)

// Payload 客户端给我们发来的数据(解析完的)
type Payload struct {
	Data  resp.Reply
	Error error
}

// 解析器parser的状态
type readState struct {
	readingMultiLine bool // 解析器正在解析单行数据还是多行数据
	// 正在读取的指令有几个参数
	// eg： set key value 三个参数 expectedArgsCount=3
	expectedArgsCount int
	msgType           byte     // 记录的消息的类型
	args              [][]byte // 用户传过来的具体的数据本身
	bulkLen           int64    // 指令的长度
}

func (r *readState) finished() bool {
	return r.expectedArgsCount > 0 && len(r.args) == r.expectedArgsCount
}

// ParseStream 异步进行解析指令，每个用户有一个解析器，会开启一个协程
func ParseStream(reader io.Reader) <-chan *Payload {

	// 异步进行解析，使得执行命令的同时还可以解析指令
	ch := make(chan *Payload)
	go parse0(reader, ch)
	return ch
}

func parse0(reader io.Reader, ch chan<- *Payload) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(debug.Stack())
		}
	}()
	bufReader := bufio.NewReader(reader)

	var (
		state readState
		err   error
		msg   []byte
	)

	for {
		var ioErr bool
		msg, ioErr, err = readLine(bufReader, &state)
		if err != nil {
			// 有ioErr，直接返回，结束
			if ioErr {
				ch <- &Payload{
					Error: err,
				}
				close(ch)
				return
			}
			ch <- &Payload{
				Error: err,
			}
			state = readState{}
			continue
		}

		// 是不是多行解析模式，这个数据是多行的，且这行数据是否正在被读
		if !state.readingMultiLine {
			// * 表示用户输入是一个数组
			if msg[0] == '*' {
				if err := parseMultiBulkHeader(msg, &state); err != nil {
					ch <- &Payload{
						Error: errors.New("protocol error: " + string(msg)),
					}
					state = readState{}
					continue
				}
				if state.expectedArgsCount == 0 {
					ch <- &Payload{
						Data: reply.NewEmptyMultiBulkReply(),
					}
					state = readState{}
					continue
				}
			} else if msg[0] == '$' { // $ 表示用户输入是一个字符串
				if err := parseBulkHeader(msg, &state); err != nil {
					ch <- &Payload{
						Error: errors.New("protocol error: " + string(msg)),
					}
					state = readState{}
					continue
				}
				// 用户发来的是 $-1\r\n, 空指令
				if state.bulkLen == -1 {
					ch <- &Payload{
						Data: reply.NewNullBulkReply(),
					}
					state = readState{}
					continue
				}
			} else {
				// 遇见+，-，:
				result, err := parseSingleLineReply(msg)
				ch <- &Payload{
					Data:  result,
					Error: err,
				}
				state = readState{}
				continue
			}
		} else {
			if err := readBody(msg, &state); err != nil {
				ch <- &Payload{
					Error: err,
				}
				state = readState{}
				continue
			}
			// 判断这整条语句是否读取完成
			if state.finished() {
				var result resp.Reply
				switch state.msgType {
				case '*':
					result = reply.NewMultiBulkReply(state.args)
				case '$':
					result = reply.NewMultiBulkReply(state.args)
				}
				ch <- &Payload{
					Data: result,
				}
				state = readState{}
			}
		}
	}
}

// 读取一个指令
// *3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n
// 一行的一行读取
func readLine(bufReader *bufio.Reader, state *readState) ([]byte, bool, error) {
	var (
		msg []byte
		err error
	)

	// 1.按照\r\n进行切分（有问题，传来的数据其中可能有\r\n）
	if state.bulkLen == 0 {
		msg, err = bufReader.ReadBytes('\n')
		// io错误
		if err != nil {
			return nil, true, err
		}
		// 其他错误
		if len(msg) == 0 || msg[len(msg)-2] != '\r' {
			return nil, false, errors.New("protocol error: " + string(msg))
		}
	} else {
		// 2.之前读到了$数字，严格按照字符个数读取
		// 实际的指令，加上\r\n
		msg = make([]byte, state.bulkLen+2)
		_, err = io.ReadFull(bufReader, msg)
		if err != nil {
			return nil, true, err
		}

		if len(msg) == 0 || msg[len(msg)-2] != '\r' || msg[len(msg)-1] != '\n' {
			return nil, false, errors.New("protocol error: " + string(msg))
		}
		state.bulkLen = 0
	}

	return msg, false, nil
}

// 读*3\r\n的内容，改变解析器的状态
func parseMultiBulkHeader(msg []byte, state *readState) error {
	var (
		err          error
		expectedLine uint64
	)

	expectedLine, err = strconv.ParseUint(string(msg[1:len(msg)-2]), 10, 64)
	if err != nil {
		return errors.New("protocol error: " + string(msg))
	}

	if expectedLine == 0 {
		state.expectedArgsCount = 0
		return nil
	} else if expectedLine > 0 {
		state.msgType = msg[0]
		state.readingMultiLine = true
		state.expectedArgsCount = int(expectedLine)
		state.args = make([][]byte, 0, expectedLine)
		return nil
	}
	return errors.New("protocol error: " + string(msg))
}

// $4\r\nPING\r\n
func parseBulkHeader(msg []byte, state *readState) error {
	var err error
	state.bulkLen, err = strconv.ParseInt(string(msg[1:len(msg)-2]), 10, 64)
	if err != nil {
		return errors.New("protocol error: " + string(msg))
	}
	if state.bulkLen == -1 {
		return nil
	} else if state.bulkLen > 0 {
		state.msgType = msg[0]
		state.readingMultiLine = true
		state.expectedArgsCount = 1
		state.args = make([][]byte, 0, 1)
		return nil
	} else {
		return errors.New("protocol error: " + string(msg))
	}
}

// 客户端发送+OK -Err或者数字
func parseSingleLineReply(msg []byte) (resp.Reply, error) {
	str := strings.TrimSuffix(string(msg), "\r\n")
	var result resp.Reply
	switch msg[0] {
	case '+':
		result = reply.NewStatusReply(str[1:])
	case '-':
		result = reply.NewErrReply(str[1:])
	case ':':
		number, err := strconv.ParseInt(str[1:], 10, 64)
		if err != nil {
			return nil, errors.New("protocol error: " + string(msg))
		}
		reply.NewIntReply(number)
	}
	return result, nil
}

// $3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n  读取之后的内容
func readBody(msg []byte, state *readState) error {
	// 去除最后面的\r\n
	line := msg[:len(msg)-2]
	var err error

	// $3
	if line[0] == '$' {
		state.bulkLen, err = strconv.ParseInt(string(line[1:]), 10, 64)
		if err != nil {
			return errors.New("protocol error: " + string(msg))
		}
		if state.bulkLen <= 0 {
			state.args = append(state.args, []byte{})
			state.bulkLen = 0
		}
	} else {
		state.args = append(state.args, line)
	}

	return nil
}
