### 示例



#### 测试服务

```
// 用于测试的结构体
type UserService struct {
   Id   int
   Name string
   Age int
}

// 用于测试查询用户的方法
func (u *UserService)QueryUser(uid int) (*UserService, bool) {
   users := make([]UserService,10)
   for i := 0; i < len(users); i++ {
      users[i]=UserService{
         Id:   i,
         Name: "用户"+strconv.Itoa(i),
         Age:  20+i,
      }
   }

   for _, user := range users {
      if user.Id==uid{
         return &user,true
      }
   }

   return &UserService{}, false
}

func (u *UserService) Test() string {
   return "hello"
}
```



#### 服务端

```
addr := "127.0.0.1:8080"
// 创建服务端
//中间键可不加
mid:= transfer.NewServiceMiddleware(func(ctx context.Context, params ...interface{}) error {
   if params[0]=="hello"{
      return nil
   }
   return transfer.ERRPRIVILEGE
})
srv := NewServer(addr,*mid)
// 将方法注册到服务端
err := srv.Register(&model.UserService{})
if err!=nil{
   fmt.Println("register service err:",err)
}
// 服务端等待调用
 srv.Run()
```

### 客户端

```
addr := "127.0.0.1:8080"
conn , err := net.Dial("tcp", addr)
if err != nil {
   t.Error(err)
}
// 创建客户端
cli := NewClient(conn)
mid:= transfer.NewClientMiddleware("hello")
err = cli.Register(new(model.UserService),*mid)
if err!=nil{
   log.Println(err)
   return
}
//添加超时
ctx,_:=context.WithTimeout(context.Background(),5*time.Second)
// 调用方法
res,err := cli.CallFunc(ctx,0, 1)
if err!=nil{
   log.Println(err)
   return
}
fmt.Println(res[0])


//通过CallFunc调用方法，上下文可添加超时，num为服务的方法按首字母排序的序号，params参数。
func (c *Client) CallFunc(ctx context.Context, Num int, params ...interface{}) ([]interface{},error) 

```

