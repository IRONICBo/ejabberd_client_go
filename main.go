package main

import (
	"bufio"
	"ejabberd_client_go/utils"
	"ejabberd_client_go/xmppc"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/shirou/gopsutil/load"
	"github.com/sirupsen/logrus"
)

var (
	// counter utils.Counter  // 计数器
	wg sync.WaitGroup // 用于等待所有的goroutine结束
)

var counter = utils.NewCounter()

// 询问一下，这个counter是不是全局变量，还是局部变量，还是说是一个指针？
// 为什么这里的counter是一个指针，而不是一个变量？
// 为什么counter在全局变量中，main函数可以获得值
// 但是在main中初始化，在goroutine中就不能获得值了？
func main() {
	// 将日志输出到文件
	file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logrus.SetOutput(file)
	} else {
		logrus.Info("Failed to log to file, using default stderr")
	}

	// 命令行参数
	Name := flag.String("name", "admin", "username")
	Password := flag.String("password", "password", "password")
	flag.Parse()

	logrus.Info(*Name, *Password)

	// 建立连接
	xc := xmppc.NewXmppConn(*Name, *Password)
	xc.Connect()

	// 获取命令行输入
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Input info with name and msg: ")
		text, _ := reader.ReadString('\n')
		infos := strings.Split(text, " ")

		if len(infos) != 2 {
			logrus.Error("Error input")
			continue
		}

		xc.SendMsg(infos[1], infos[0])
	}

	info, _ := load.Avg()
	fmt.Printf("%v\n", info)

	fmt.Printf("Pointer %v\n", &counter)
	benchmark()

	// 打印goroutine数量
	wg.Wait()
	fmt.Println("Done...")
	fmt.Println("goroutine num: ", runtime.NumGoroutine())
	fmt.Println("Total Msg Count Num: ", counter.Value())
}

func benchmark() {
	for i := 1; i <= 2; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			xc := xmppc.NewXmppConn(fmt.Sprintf("light%03d", i), fmt.Sprintf("light%03d", i))
			// 打印连接状态
			connectStatus := xc.Connect()
			logrus.Infof("Connect status: %s", connectStatus)

			for j := 0; j < 1; j++ {
				// 发送的消息内容是当前时间和发送用户的名字以及消息序号
				// xc.SendMsg(fmt.Sprintf("%s %s %d", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf("light%03d", i), i), fmt.Sprintf("light%03d", i))
				// xc.SendMsg(fmt.Sprintf("%s %s %d", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf("light%03d", i), i), fmt.Sprintf("light%03d", (i+100)%200))
				xc.SendMsg(fmt.Sprintf("%s %s %d", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf("light%03d", i), i), fmt.Sprintf("light%03d", i+100))

				counter.Inc()

				logrus.Infof("Send msg to %s: %s", fmt.Sprintf("light%03d", i), fmt.Sprintf("%s %s %d", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf("light%03d", i), i))

				// runtime打印系统cpu，内存占用
				// v, _ := mem.VirtualMemory()
				// fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)
				// // fmt.Println(v)
				// info, _ := load.Avg()
				// fmt.Printf("%v\n", info)

				time.Sleep(time.Second)
			}
		}(i)
	}
}

// unknown namespace urn:ietf:params:xml:ns:xmpp-streams <host-unknown/>
// 这里的地址有变化

// [warning] (tcp|<0.649.0>) Failed to secure c2s connection: TLS failed: no_certfile
// ERRO[0000] Xmpp connect failed: expecting starttls proceed: expected element type <proceed> but have <failure>
// 连接失败了，没有证书
