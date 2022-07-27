/*******
* @Author:qingmeng
* @Description:
* @File:rpc_test
* @Date:2022/7/24
 */

package main

import (
	"context"
	"dur-rpc/model"
	"dur-rpc/transfer"
	"fmt"
	"log"
	"net"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	// 需要对interface可能产生的类型进行注册
	addr := "127.0.0.1:8080"
	// 创建服务端
	//mid := transfer.NewServiceMiddleware(func(ctx context.Context, params ...interface{}) error {
	//	if params[0] == "hello" {
	//		return nil
	//	}
	//	return transfer.ERRPRIVILEGE
	//})
	srv := NewServer(addr)
	// 将方法注册到服务端
	err := srv.Register(&model.UserService{})
	if err != nil {
		fmt.Println("register service err:", err)
	}
	// 服务端等待调用
	srv.Run()

}

func TestClient(t *testing.T) {
	// 客户端获取连接
	addr := "127.0.0.1:8080"
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Error(err)
	}
	// 创建客户端
	cli := NewClient(conn)
	err = cli.Register(new(model.UserService))
	if err != nil {
		log.Println(err)
		return
	}
	//添加超时
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	// 调用方法
	res, err := cli.CallFunc(ctx, 0, 1)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(res[0])

	go func() {
		res, err = cli.CallFunc(ctx, 1)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(res)
	}()

	for i := 0; i < 5; i++ {
		time.Sleep(time.Second)
		go func() {
			for i := 0; i < 5; i++ {
				res, err = cli.CallFunc(ctx, 0, i)
				if err != nil {
					log.Println(err)
					return
				}
				fmt.Println(res[0], res[1])
			}

		}()
	}
}

func TestClient1(t *testing.T) {
	// 客户端获取连接
	addr := "127.0.0.1:8080"
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Error(err)
	}
	// 创建客户端
	mid := transfer.NewClientMiddleware("hello")
	cli := NewClient(conn)
	err = cli.Register(new(model.UserService), *mid)
	if err != nil {
		log.Println("register err", err)
		return
	}
	// 调用方法
	res, err := cli.CallFunc(context.Background(), 0, 1)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(res[0])
	for i := 0; i < 7; i++ {
		res, err = cli.CallFunc(context.Background(), 0, i)
		if err != nil {
			log.Println(err)
			return
		}
		time.Sleep(1 * time.Second)
		fmt.Println(res[0], res[1])

	}

}
