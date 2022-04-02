Go 实现的轻量级的事件管理、调度工具库

- 支持自定义事件对象
- 支持同时监听多个事件


## 主要方法


- `Subscribe(listener Listener) error`  订阅，支持注册多个事件监听
- `Publish(event interface{}) ` 触发事件


## 快速使用

```go
package main
import (
	"fmt"
	
	"github.com/lianglong/go/event"
)

type user struct{
    Id int
    Name string
}

type UserRegistered struct{
    Data *user
}

type ChangeUserPassWord struct{
    Data *user
}

type SendEmailListener struct{
}

func (listener SendEmailListener)Listen() []interface{} {
    return []interface{}{
        UserRegistered{},
        ChangeUserPassWord{},
    }
}

func (listener SendEmailListener)Handle(e interface{}) error {
    switch ev := e.(type) {
    case UserRegistered:
        fmt.Printf("send email for registered user:%s",ev.Data.Name)
    case ChangeUserPassWord:
        fmt.Printf("send email for change password user:%s",ev.Data.Name)
    }
    return nil
}

func (listener SendEmailListener)Priority() int  {
    return event.NormalPriority
}

func main() {
    eventDispatcher := event.New()
	// 注册事件监听器
	eventDispatcher.Subscribe(SendEmailListener{})
	
	// ... ...
	
	// 触发事件
	eventDispatcher.Publish(UserRegistered{&user{Id:1,Name:"user1"}})
}
```


## LICENSE

**[MIT](LICENSE)**