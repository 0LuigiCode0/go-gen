package tmp

const HandlerMQTTTmp = `package {{printf "%v_handler" (index . 0)}}

import (
	"encoding/json"
	"fmt"
	"strings"
	{{printf "\"%v/handlers/%v_handler/%v_helper\"" (index . 1) (index . 0) (index . 0)}}
	{{printf "\"%v/helper\"" (index . 1)}}
	{{printf "\"%v/hub/hub_helper\"" (index . 1)}}
	
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type handler struct {
	hub_helper.HelperForHandler

	client mqtt.Client
	err    chan error
}

func InitHandler(hub hub_helper.HelperForHandler, conf *helper.HandlerConfig) (H {{print (index . 0)}}_helper.Handler, err error) {
	h := &handler{HelperForHandler: hub}
	H = h

	client := mqtt.NewClient(mqtt.NewClientOptions().
		AddBroker(fmt.Sprintf("tcp://%v:%v", conf.Host, conf.Port)).
		SetUsername(conf.User).
		SetPassword(conf.Password).
		SetOnConnectHandler(h.onConn).
		SetConnectionLostHandler(h.onLost).
		SetConnectRetry(true).
		SetResumeSubs(true))
	h.client = client
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	helper.Wg.Add(1)
	go h.loop()

	return
}

func (h *handler) haron(_ mqtt.Client, msg mqtt.Message) {
	var topic, route string
	str := strings.Split(msg.Topic(), "/")
	if len(str) == 3 {
		route = strings.TrimSpace(str[1])
		topic = strings.TrimSpace(str[2])
	} else {
		return
	}
	switch route {
	}
}

func (h *handler) onConn(client mqtt.Client) {
	helper.Log.Service("mqtt is conected")
	if token := client.Subscribe("back/#", 0, h.haron); token.Wait() && token.Error() != nil {
		h.err <- token.Error()
	}
}

func (h *handler) onLost(client mqtt.Client, err error) {
	helper.Log.Servicef("mqtt is disconected: %v", err)
	token := client.Unsubscribe("back/#")
	token.Wait()
}

func (h *handler) loop() {
	defer helper.Wg.Done()

	for {
		select {
		case e := <-h.err:
			helper.Log.Errorf("mqtt erorr: %v", e)
			h.client.Disconnect(1000)
			return
		case <-helper.Ctx.Done():
			helper.Log.Service("mqtt stoped")
			h.client.Disconnect(1000)
			return
		}
	}
}

func (h *handler) respOk(topic string, data interface{}) {
	resp := &helper.ResponseModel{
		Success: true,
		Result:  data,
	}
	buf, err := json.Marshal(resp)
	if err != nil {
		helper.Log.Warningf(helper.KeyErrorParse+": json: %v", err)
		return
	}
	if token := h.client.Publish(topic, 0, false, buf); token.Wait() && token.Error() != nil {
		helper.Log.Warningf(helper.KeyErrorSend+": mqtt: %v", token.Error())
	}
}

func (h *handler) respError(topic string, code helper.ErrCode, msg string) {
	resp := &helper.ResponseModel{
		Success: false,
		Result: &helper.ResponseError{
			Code: code,
			Msg:  msg,
		},
	}
	buf, _ := json.Marshal(resp)
	if token := h.client.Publish(topic, 0, false, buf); token.Wait() && token.Error() != nil {
		helper.Log.Warningf(helper.KeyErrorSend+": mqtt: %v", token.Error())
	}
}`
