package mqtt
import "github.com/eclipse/paho.mqtt.golang"

func init() {
    client:=mqtt.NewClient(nil)
    client.Connect()
}

