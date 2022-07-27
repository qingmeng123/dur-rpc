/*******
* @Author:qingmeng
* @Description:
* @File:session
* @Date:2022/7/24
 */

package transfer

import (
	"encoding/binary"
	"log"
	"net"
)


// 会话连接的结构体
type Session struct {
	conn net.Conn
}
// 创建新连接
func NewSession(conn net.Conn) *Session {
	return &Session{conn: conn}
}
// 向连接中写数据
func (s Session) Write(data []byte) error {
	//先写入头部数据，记录数据长度,防止丢包
	buf := make([]byte, 4)
	// binary 只认固定长度的类型，所以使用了uint32，而不是直接写入
	binary.BigEndian.PutUint32(buf[:4], uint32(len(data)))
	n, err := s.conn.Write(buf)
	if n!=4||err!=nil{
		return err
	}

	//再发送data
	n,err=s.conn.Write(data)
	if n!=len(data)||err!=nil{
		return err
	}
	return nil
}

// 从连接中读数据
func (s Session) Read() ([]byte, error) {
	// 读取头部长度
	header := make([]byte, 4)
	_, err := s.conn.Read(header)
	if err != nil {
		return nil, err
	}
	dataLen := binary.BigEndian.Uint32(header)
	// 按照头部长度作为最长去读取数据
	data := make([]byte, dataLen)
	n, err := s.conn.Read(data)
	if n!=int(dataLen)|| err != nil {
		return nil, err
	}
	return data, nil
}

//将数据解码获取
func (s Session) GetData() (rpcData RPCData,err error) {
	b, err := s.Read()
	if err != nil {
		return
	}
	// 对数据解码
	err = rpcData.Decode(b)
	if err != nil {
		return
	}
	return
}


//将数据编码发送
func (s Session) SendData(data RPCData)  {
	// 编码
	respBytes, err := data.Encode()
	if err != nil {
		data.Args=nil
		data.Err= ERRENCODE.Error()
		s.SendData(data)
		return
	}

	err = s.Write(respBytes)
	if err != nil {
		return
	}
	log.Println("发送回复",data)
	return
}
