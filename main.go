package main

import (
	"flag"
	"fmt"
	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pool/goroutine"
	"gorm.io/gorm"
	"log"
	"log-transfer/core"
	"log-transfer/protocol"
	"time"
)

/**
@author zhangchengji
@time 2021-03-08 竣工
*/
var (
	NewMysql *core.Mysql
	logList  chan string
	DB       *gorm.DB
)

type logCodecServer struct {
	*gnet.EventServer
	addr       string
	multicore  bool
	async      bool
	codec      gnet.ICodec
	workerPool *goroutine.Pool
}

func (cs *logCodecServer) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	log.Printf("启动成功🚀 监听端口 %s (multi-cores: %t, loops: %d)\n",
		srv.Addr.String(), srv.Multicore, srv.NumEventLoop)
	return
}

func (cs *logCodecServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	//TODO 这么没有加连接关闭是因为，客户端 会自动断连，无需担心
	logList <- string(frame)

	return
}
func (es *logCodecServer) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	log.Println("长时间没有日志传输，自动断连...🛵")
	return
}

/**
多路复用读取
*/
func read() {

	for {
		select {
		case log := <-logList:
			core.SendChanLog(log)

		default:
			time.Sleep(time.Microsecond * 100)

		}
	}

}
func logCodecServe(addr string, multicore, async bool, codec gnet.ICodec) {
	var err error
	codec = &protocol.LogLengthFieldProtocol{}
	cs := &logCodecServer{addr: addr, multicore: multicore, async: async, codec: codec, workerPool: goroutine.Default()}
	err = gnet.Serve(cs, addr, gnet.WithMulticore(multicore), gnet.WithTCPKeepAlive(time.Minute*5), gnet.WithCodec(codec))
	if err != nil {
		panic(err)
	}
}
func init() {
	NewMysql = &core.Mysql{
		Username:     "root",
		Password:     "123456",
		Path:         "localhost",
		Port:         3306,
		Dbname:       "log",
		Config:       "charset=utf8mb4&parseTime=True&loc=Local",
		MaxIdleConns: 10,
		MaxOpenConns: 100,
		LogMode:      false,
	}
}
func main() {
	/**
	gorm 配置加载
	*/
	DB = core.GormMysql(NewMysql)
	core.MysqlTables(DB)
	dvb, _ := DB.DB()
	defer dvb.Close()
	/**
	启动日志收集处理
	*/
	logList = make(chan string, 1000)

	go read()

	core.Init()
	var port int
	var multicore, reuseport bool

	// Example command: go run server.go --port 9000 --multicore=true
	flag.IntVar(&port, "port", 9000, "server port")
	flag.BoolVar(&multicore, "multicore", true, "multicore") //多核处理器全部使用
	flag.BoolVar(&reuseport, "reuseport", false, "--reuseport true")

	flag.Parse()
	addr := fmt.Sprintf("tcp://:%d", port)
	logCodecServe(addr, multicore, false, nil)

}
