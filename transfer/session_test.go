/*******
* @Author:qingmeng
* @Description:
* @File:session_test
* @Date:2022/7/24
 */

package transfer

import (
	"fmt"
	"net"
	"sync"
	"testing"
)

func TestSession(t *testing.T) {
	addr := "127.0.0.1:8080"
	hello := "hello world"
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		// 创建tcp连接
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			t.Fatal(err)
		}
		conn,_ := lis.Accept()
		s := Session{conn: conn}
		err = s.Write([]byte(hello))
		if err != nil {
			t.Fatal(err)
		}
	}()
	go func() {
		defer wg.Done()
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			t.Fatal(err)
		}
		s := Session{conn: conn}
		// 读数据
		data, err := s.Read()
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != hello {
			t.Fatal(err)
		}
		fmt.Println(string(data))
	}()
	wg.Wait()
}