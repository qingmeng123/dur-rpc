/*******
* @Author:qingmeng
* @Description:
* @File:client
* @Date:2022/7/24
 */

package main

import (
	"context"
	"dur-rpc/transfer"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"reflect"
	"time"
)

// 声明客户端
type Client struct {
	// 请求的服务名
	service reflect.Value
	conn net.Conn
	mid  transfer.Middleware
}

var defaultMiddleware = transfer.Middleware{
	On:          false,
	Interceptor: nil,
	Params: nil,
}

// 创建客户端对象
func NewClient(conn net.Conn) *Client {
	return &Client{ conn: conn}
}

//注册服务
func (c *Client) Register(req interface{},mids...transfer.Middleware)error  {
	gob.Register(req)
	c.service=reflect.ValueOf(req)
	mid:= defaultMiddleware
	if len(mids)>0{
		gob.Register(mids[0])
		mid=mids[0]
	}
	c.mid=mid
	return c.Dial()
}

//处理请求
func (c *Client) handleReq(reqRPC transfer.RPCData) transfer.RPCData {
	cliSession := transfer.NewSession(c.conn)
	//发送请求
	b, err := reqRPC.Encode()
	if err != nil {
		panic(err)
	}
	err = cliSession.Write(b)
	if err != nil {
		panic(err)
	}

	//接收回复
	respBytes, err := cliSession.Read()
	if err != nil {
		panic(err)
	}
	err = reqRPC.Decode(respBytes)
	if err != nil {
		panic(err)
	}

	return reqRPC
}

//尝试连接
func (c *Client) Dial() error {
	resp := c.handleReq(transfer.RPCData{Mid: c.mid})
	fmt.Println("resp",resp)
	if resp.Err!=""{
		return transfer.ERRPRIVILEGE
	}
	return nil
}

//调用方法
func (c *Client) CallFunc(ctx context.Context, Num int, params ...interface{}) ([]interface{},error) {

	respChan := make(chan transfer.RPCData, 5)
	go func() {
		respChan <- c.handleReq(transfer.RPCData{Num: Num, Args: params})
	}()
	respData := transfer.RPCData{}

		select {
		case respData = <-respChan:
			if respData.Err != "" {
				log.Println(respData.Err)
				return nil,nil
			}
			return respData.Args,nil
			//用户自定义超时时间
		case <-ctx.Done():
			<-respChan
			return nil, transfer.ERRTIME
			//默认3s超时
		case <-time.After(3 * time.Second):
			<-respChan
			return nil, transfer.ERRTIME
	}
}


