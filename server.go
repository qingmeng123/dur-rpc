/*******
* @Author:qingmeng
* @Description:
* @File:server
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
	"sync"
)

// 声明服务端
type Server struct {
	// 地址
	addr string
	// 服务端维护的服务名
	service reflect.Value
	mid     transfer.Middleware
}

// 创建服务端对象
func NewServer(addr string, mids ...transfer.Middleware) *Server {
	mid := defaultMiddleware
	if len(mids) > 0 {
		mid = mids[0]
	}
	return &Server{addr: addr, mid: mid, service: reflect.Value{}}
}

//服务端注册服务
func (s *Server) Register(f interface{}) error {
	gob.Register(f)
	gob.Register(s.mid)
	val := reflect.ValueOf(f)
	fmt.Println("注册服务成功,方法个数：", val.Type().NumMethod())
	//传参结构体错误
	if val.Kind() != reflect.Ptr {
		return transfer.ERRTYPE
	}
	elem := val.Type().Elem()
	if elem.Kind() != reflect.Struct {
		return transfer.ERRTYPE
	}

	//将结构体指针的映射添加到服务中
	s.service = val
	return nil
}

//处理请求
func (s *Server) handleRequest(reqData transfer.RPCData) (respData transfer.RPCData) {
	//验证中间键
	fmt.Println("req", reqData)
	if s.mid.On {
		//验证请求是否添加中间键
		if !reqData.Mid.On {
			respData.Err = transfer.ERRPRIVILEGE.Error()
			return respData
		}
		err := s.mid.Interceptor(context.Background(), reqData.Mid.Params...)
		if err != nil {
			log.Println(err)
			respData.Err = transfer.ERRPRIVILEGE.Error()
			return reqData
		}
		s.mid.On = false
		return respData
	}

	//判断请求是否正确
	numOfMethod := s.service.NumMethod()
	if reqData.Num >= numOfMethod {
		log.Println(transfer.ERRNUM)
		respData.Err = transfer.ERRNUM.Error()
		return
	}
	method := s.service.Method(reqData.Num)
	if method.Type().NumIn() != len(reqData.Args) {
		log.Println(transfer.ERRPARAMS)
		respData.Err = transfer.ERRPARAMS.Error()
		return
	}

	// 解析遍历客户端出来的参数, 放到一个数组中
	inArgs := make([]reflect.Value, 0, len(reqData.Args))
	for _, arg := range reqData.Args {
		inArgs = append(inArgs, reflect.ValueOf(arg))
	}
	// 反射调用方法，传入参数，起一个协程处理方法，避免方法本身panic
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				respData.Num = reqData.Num
				respData.Err = transfer.ERRSERVICEFUNC.Error()
			}
		}()
		resp := method.Call(inArgs)
		// 解析遍历执行结果，放到一个数组中
		outArgs := make([]interface{}, 0, len(resp))
		for _, o := range resp {
			outArgs = append(outArgs, o.Interface())
		}
		respData.Num = reqData.Num
		respData.Args = outArgs
	}()
	wg.Wait()
	return
}

//处理连接
func (s *Server) serverConn(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println("close err:", err)
		}
	}(conn)
	wg := new(sync.WaitGroup)
	for {
		srvSession := transfer.NewSession(conn)
		//创建回复
		var respData transfer.RPCData
		// RPC 读取数据
		reqData, err := srvSession.GetData()
		if err != nil {
			respData.Err = transfer.ERRDATA.Error()
			srvSession.SendData(respData)
			return
		}
		wg.Add(1)
		//处理请求
		go func() {
			defer wg.Done()
			respData = s.handleRequest(reqData)
			//发送回复
			//time.Sleep(6*time.Second)//模拟超时
			srvSession.SendData(respData)
			return
		}()
		wg.Wait()
	}
}

// Run 开始监听
func (s *Server) Run() {
	// 监听
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		fmt.Printf("监听%s err:%v", s.addr, err)
		return
	}
	//记录是否开启拦截器
	state := s.mid.On

	for {
		// 拿到连接
		conn, err := lis.Accept()
		if err != nil {
			fmt.Printf("accept err:%v", err)
			return
		}
		// 处理连接
		s.mid.On = state
		go s.serverConn(conn)
	}
}
