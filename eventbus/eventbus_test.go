package eventbus

import (
	"fmt"
	"testing"
)

func TestEvent(t *testing.T) {
	message := EventBus.CreateMessage()
	message.Set("testkey", "testvalue")

	client := EventBus.CreateClient(func(message *Message) {
		fmt.Print(message.Get("testkey"))
	})

	EventBus.On("test", client)
	EventBus.Emit("test", message)
}
