/*******
* @Author:qingmeng
* @Description:
* @File:middleware
* @Date:2022/7/27
 */

package transfer

import (
	"bytes"
	"context"
	"encoding/gob"
)

//中间键
type Middleware struct {
	On          bool             //是否启用中间件
	Interceptor UnaryInterceptor //拦截器
	Params      []interface{}    //给拦截器用的参数
}

//服务端中间键
func NewServiceMiddleware(interceptor UnaryInterceptor) *Middleware {
	return &Middleware{
		On:          true,
		Interceptor: interceptor,
	}
}

//客户端中间键
func NewClientMiddleware(params ...interface{}) *Middleware {
	return &Middleware{
		On:       true,
		Params: params,
	}
}

func (m *Middleware) Encode() ([]byte, error) {
	var buf bytes.Buffer
	// 得到字节数组的编码器
	bufEnc := gob.NewEncoder(&buf)
	// 对数据进行编码
	if err := bufEnc.Encode(m); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (m *Middleware) Decode(b []byte) error {
	buf := bytes.NewBuffer(b)
	// 返回字节数组的解码器
	bufDec := gob.NewDecoder(buf)
	// 对数据解码
	if err := bufDec.Decode(&m); err != nil {
		return  err
	}
	return nil
}

type UnaryInterceptor func(ctx context.Context,params ...interface{})error
