package queue

import (
	"errors"
	"fmt"
	"github.com/golang-module/carbon"
	"owl"
	"owl/contract"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

var (
	connections = make(map[string]*amqp.Connection, 10)
	linkMap     = make(map[string]*link, 10)
	lock        sync.Mutex
)

type link struct {
	dsn string
}

type Options struct {
	contract.ServerConfig
	Vhost   string `json:"vhost"`
	CfgFile string `json:"cfg-file"`
	AbsPath string `json:"abs-path"`
}

func NewOption(cfgManager *owl.ConfManager) (opt *Options) {
	err := cfgManager.GetConfig("rabbit", &opt)
	if err != nil {
		return nil
	}
	opt.CfgFile = "rabbit"
	return opt
}

type RabbitMQ struct {
	opt          *Options
	l            contract.Logger
	clientNotify chan *amqp.Error
	exchange     string
	queue        string
	routingKey   string
	closeNotify  chan *amqp.Error
}

func NewRabbit(opt *Options, l *owl.LoggerFactory) *RabbitMQ {
	r := &RabbitMQ{
		opt: opt,
		l:   l.RuntimeLogger(),
	}
	return r
}

// Connect  多个实例公共使用一个连接
func (i *RabbitMQ) Connect() *amqp.Connection {
	var err error
	lock.Lock()
	defer lock.Unlock()

	dbLinkInfo, ok := linkMap[i.opt.CfgFile] // 申请连接过
	if !ok {

		dsn := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
			i.opt.Username,
			i.opt.Password,
			i.opt.Host,
			i.opt.Port,
			i.opt.Vhost,
		)
		linkMap[i.opt.CfgFile] = &link{dsn: dsn}
		dbLinkInfo = linkMap[i.opt.CfgFile]
		go func() {
			for {
				select {
				case file := <-owl.CfgChangeNotify[i.opt.AbsPath]:
					lock.Lock()
					i.l.Info("重载配置" + file)
					delete(linkMap, i.opt.CfgFile)
					delete(connections, dsn)
					lock.Unlock()
					return
				}
			}
		}()
	}

	dsn := dbLinkInfo.dsn

	con, ok := connections[dsn] // 申请连接过

	if !ok || con == nil || (con != nil && con.IsClosed()) {
		con, err = amqp.Dial(dsn) // 尝试连接一次
		connections[dsn] = con

		if err != nil {
			return nil
		}

		return con
	} else {
		return con
	}
}

func (i *RabbitMQ) Queue(queue string) *RabbitMQ {
	i.queue = queue
	return i
}
func (i *RabbitMQ) newQueue() error {
	con := i.Connect()
	if con == nil {
		return errors.New("rabbit 未连接")
	}
	ch, err := con.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	_, err = ch.QueueDeclare(
		i.queue, // 队列名称
		false,   // 持久性
		false,   // 自动删除
		false,   // 排他性
		false,   // 不等待
		nil,     // 其他属性
	)

	return err
}

// Publish 发布消息到 RabbitMQ
func (i *RabbitMQ) Publish(message []byte) error {
	con := i.Connect()
	if con == nil {
		return errors.New("rabbit 还未连接")
	}

	err := i.newQueue()
	if err != nil {
		return err
	}

	ch, err := con.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.Publish(
		i.exchange, // 交换机名称
		i.queue,    // 路由键
		false,      // 强制
		false,      // 立即
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         message,
			DeliveryMode: amqp.Persistent, // 持久性
		},
	)
	return err
}

// Consume 使用消费者处理消息
func (i *RabbitMQ) Consume(prefetchCount int, handler func(data string, msg amqp.Delivery, i *int32)) {

	go func() {
	ReConsume:
		con := i.Connect()

		if con == nil {
			time.Sleep(time.Second * 3)
			fmt.Println("消费者重连中")
			goto ReConsume
		}

		err := i.newQueue()
		if err != nil {
			time.Sleep(time.Second * 3)
			goto ReConsume
		}

		conClose := con.NotifyClose(make(chan *amqp.Error))
		i.l.Debug("开始消费")

		i.l.Debug("获取 channle")
		ch, err := con.Channel()
		i.l.Debug("获取 channle end")
		if err != nil {
			i.l.Debug("获取通道失败")
			time.Sleep(time.Second * 3)
			goto ReConsume
		}

		err = ch.Qos(prefetchCount, 0, false)
		if err != nil {
			goto ReConsume
		}
		msgs, err := ch.Consume(
			i.queue, // 队列名称
			"",      // 消费者标签
			false,   // 自动应答
			false,   // 排他性
			false,   // 不等待
			false,   // 其他属性
			nil,
		)
		if err != nil {
			//fmt.Println("队列消费失败", err)
			goto ReConsume
		}
		var running int32
		var exitFor bool
		for {
			select {
			case _, ok := <-conClose:
				if !ok {
					i.l.Debug("通道关闭")
				}
				i.l.Debug("收到退出信号")
				exitFor = true
				break
			case msg, ok := <-msgs:
				if !ok {
					i.l.Debug("通道关闭")
					exitFor = true
					break
				}

				go handler(string(msg.Body), msg, &running)

			default:
				i.l.Debug("消费者获取数据中", carbon.Now().ToDateTimeString())
			}
			if exitFor {
				ch.Close()
				goto ReConsume
			}
		}
	}()
}
