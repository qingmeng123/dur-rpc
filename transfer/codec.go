/*******
* @Author:qingmeng
* @Description:
* @File:codec
* @Date:2022/7/24
 */

package transfer

import (
	"bytes"
	"encoding/gob"
)

type Codec interface {
	Encode()([]byte,error)
	Decode([]byte) error
}

// 定义数据格式和编解码
type RPCData struct {
	// 访问的函数序号
	Num int
	// 传参或返回值
	Args []interface{}
	//返回的错误
	Err string
	Mid Middleware
}

// 编码
func (data *RPCData)Encode() ([]byte, error) {
	var buf bytes.Buffer
	// 得到字节数组的编码器
	bufEnc := gob.NewEncoder(&buf)
	// 对数据进行编码
	if err := bufEnc.Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// 解码
func (data *RPCData)Decode(b []byte) error {
	buf := bytes.NewBuffer(b)
	// 返回字节数组的解码器
	bufDec := gob.NewDecoder(buf)
	// 对数据解码
	if err := bufDec.Decode(&data); err != nil {
		return  err
	}
	return nil
}