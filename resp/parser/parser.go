package parser

import (
	"bufio"
	"errors"
	"github.com/sunflower10086/go-redis/interface/resp"
	"github.com/sunflower10086/go-redis/resp/reply"
	"io"
	"strconv"
	"strings"
)

// Payload 客户端给我们发来的数据(解析完的)
type Payload struct {
	Data  resp.Reply
	Error error
}

// 解析器parser的状态
type readState struct {
	// 解析器正在解析单行数据还是多行数据
	readingMultiLine bool
	// 正在读取的指令有几个参数
	// eg： set key value 三个参数 expectedArgsCount=3
	expectedArgsCount int
	// 记录的消息的类型
	msgType byte
	// 用户传过来的具体的数据本身
	args [][]byte
	// 指令的长度
	bulkLen int64
}

func (r *readState) finished() bool {
	return r.expectedArgsCount > 0 && len(r.args) == r.expectedArgsCount
}

func ParseStream(reader io.Reader) <-chan *Payload {
	// 异步进行解析，使得执行命令的同时还可以解析指令
	ch := make(chan *Payload)
	go parse0(reader, ch)
	return ch
}

func parse0(reader io.Reader, ch chan<- *Payload) {

}

// 读取一个指令
// *3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n
// 这个方法只知道$之后该怎么读，但是不知道*3之后怎么读
func readLine(bufReader *bufio.Reader, state *readState) ([]byte, error) {
	var (
		msg []byte
		err error
	)
	// 1.按照\r\n进行切分（有问题，传来的数据其中可能有\r\n）
	if state.bulkLen == 0 {
		msg, err = bufReader.ReadBytes('\n')
		// io错误
		if err != nil {
			if err == io.EOF {
				return msg, nil
			}
			return nil, err
		}
		// 其他错误
		if len(msg) == 0 || msg[len(msg)-2] != '\r' {
			return nil, errors.New("protocol error: " + string(msg))
		}
	} else {
		// 2.之前读到了$数字，严格按照字符个数读取
		// 实际的指令，加上\r\n
		msg = make([]byte, state.bulkLen+2)
		_, err = io.ReadFull(bufReader, msg)
		if err != nil {
			if err == io.EOF {
				return msg, nil
			}
			return nil, err
		}

		if len(msg) == 0 || msg[len(msg)-2] != '\r' || msg[len(msg)-1] != '\n' {
			return nil, errors.New("protocol error: " + string(msg))
		}
		state.bulkLen = 0
	}

	return msg, nil
}

// 读*3之后的内容，改变解析器的状态
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

// 客户端发送+OK -Err
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