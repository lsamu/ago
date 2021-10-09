package kafka

import (
    "context"
    "fmt"
    "github.com/segmentio/kafka-go"
    "log"
    "net"
    "strconv"
    "sync"
    "time"
)

func InitKafka() {

}

func GetKafka() {

}

func NewKafka() {

}

var (
    writerPool  = make(map[string]*kafka.Writer)
    writerCount = make(map[string]int) //计数
    lock        sync.Mutex
    rw          sync.RWMutex
    readerPool  = make(map[string]*kafka.Reader)
)

type (
    KafkaOrm struct {
        conn  *kafka.Conn
        batch *kafka.Batch
    }
)

func (o *KafkaOrm) Open(topic string, partition int) {
    conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
    if err != nil {
        log.Fatal("failed to dial leader:", err)
    }
    conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
    batch := conn.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max
    o.conn = conn
    o.batch = batch
}

func (o *KafkaOrm) Writer() (err error) {
    _, err = o.conn.WriteMessages(
        kafka.Message{Value: []byte("one!")},
        kafka.Message{Value: []byte("two!")},
        kafka.Message{Value: []byte("three!")},
    )
    if err != nil {
        log.Fatal("failed to write messages:", err)
    }

    if err := o.conn.Close(); err != nil {
        log.Fatal("failed to close writer:", err)
    }
    return
}

func (o *KafkaOrm) Read() {
    b := make([]byte, 10e3) // 10KB max per message
    for {
        _, err := o.batch.Read(b)
        if err != nil {
            break
        }
        fmt.Println(string(b))
    }

    if err := o.batch.Close(); err != nil {
        log.Fatal("failed to close batch:", err)
    }

    if err := o.conn.Close(); err != nil {
        log.Fatal("failed to close connection:", err)
    }
}

func (o *KafkaOrm) CreateTopic(topic string) {
    // auto.create.topics.enable='false'  KAFKA_AUTO_CREATE_TOPICS_ENABLE='false'
    controller, err := o.conn.Controller()
    if err != nil {
        panic(err.Error())
    }
    var controllerConn *kafka.Conn
    controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
    if err != nil {
        panic(err.Error())
    }
    defer controllerConn.Close()

    topicConfigs := []kafka.TopicConfig{
        kafka.TopicConfig{
            Topic:             topic,
            NumPartitions:     1,
            ReplicationFactor: 1,
        },
    }

    err = controllerConn.CreateTopics(topicConfigs...)
    if err != nil {
        panic(err.Error())
    }
}

func (o *KafkaOrm) ListTopic() {
    partitions, err := o.conn.ReadPartitions()
    if err != nil {
        panic(err.Error())
    }

    m := map[string]struct{}{}

    for _, p := range partitions {
        m[p.Topic] = struct{}{}
    }
    for k := range m {
        fmt.Println(k)
    }
}

func (o *KafkaOrm) Close() {
    defer o.conn.Close()
}
