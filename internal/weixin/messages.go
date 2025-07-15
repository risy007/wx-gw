package weixin

import (
	"encoding/json"
	"fmt"
	"github.com/xen0n/go-workwx/v2"
	"time"
)

type PrtgAlertMessage struct {
	EventTime     Time   `json:"event_time"`
	Action        string `json:"action"`
	ProbeDevice   string `json:"probe_device"`
	DeviceGroup   string `json:"device_group"`
	Node          string `json:"node"`
	SensorName    string `json:"sensor_name"`
	CurrentStatus string `json:"current_status"`
	DownTime      string `json:"down_time"`
	AttachedMsg   string `json:"attached_msg"`
	SensorID      string `json:"sensor_id"`
}

func (m *PrtgAlertMessage) ConvertToWorkWxTextMsg() *WecomTextMessage {
	body := WecomTextBody{
		Content:  m.GenWorkWxMessage("text"),
		Mentions: nil,
	}
	wtm := WecomTextMessage{
		MsgType: "text",
		Text:    body,
	}
	return &wtm
}

func (m *PrtgAlertMessage) ConvertToWorkWxMDMsg() *WecomMarkdownMessage {
	body := WecomMarkdownBody{
		Content: m.GenWorkWxMessage("md"),
	}
	wmm := WecomMarkdownMessage{
		MsgType:  "markdown",
		MarkDown: body,
	}
	return &wmm
}

func (m *PrtgAlertMessage) GenWorkWxMessage(t string) string {
	switch t {
	case "text":
		s := fmt.Sprintf("【告警消息】：%s\n", m.Action)
		s += fmt.Sprintf("【告警时间】：%s\n", m.EventTime)
		s += fmt.Sprintf("【探针设备】：%s\n", m.ProbeDevice)
		s += fmt.Sprintf("【设备分组】：%s\n", m.DeviceGroup)
		s += fmt.Sprintf("【设备节点名】：%s\n", m.Node)
		s += fmt.Sprintf("【探针名】：%s\n", m.SensorName)
		s += fmt.Sprintf("【当前状态】：%s\n", m.CurrentStatus)
		s += fmt.Sprintf("【中断时间】：%s\n", m.DownTime)
		s += fmt.Sprintf("【附加消息】：%s\n", m.AttachedMsg)
		s += fmt.Sprintf("【探针编号】：%s\n", m.SensorID)
		return s
	case "md":
		s := fmt.Sprintf("<font color=\"warning\">【告警消息】：%s</font>\n", m.Action)
		s += fmt.Sprintf("【告警时间】：<font color=\"comment\">%s</font>\n", m.EventTime)
		s += fmt.Sprintf("【探针设备】：<font color=\"comment\">%s</font>\n", m.ProbeDevice)
		s += fmt.Sprintf("【设备分组】：<font color=\"comment\">%s</font>\n", m.DeviceGroup)
		s += fmt.Sprintf("【设备节点名】：<font color=\"comment\">%s</font>\n", m.Node)
		s += fmt.Sprintf("【探针名】：<font color=\"comment\">%s</font>\n", m.SensorName)
		s += fmt.Sprintf("【当前状态】：<font color=\"comment\">%s</font>\n", m.CurrentStatus)
		s += fmt.Sprintf("【中断时间】：<font color=\"comment\">%s</font>\n", m.DownTime)
		s += fmt.Sprintf("【附加消息】：<font color=\"comment\">%s</font>\n", m.AttachedMsg)
		s += fmt.Sprintf("【探针编号】：<font color=\"comment\">%s</font>\n", m.SensorID)
		return s
	}
	js, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(js)
}

type WecomTextMessage struct {
	MsgType string        `json:"msgtype"`
	Text    WecomTextBody `json:"text"`
}

type WecomTextBody struct {
	Content  string           `json:"content"`
	Mentions *workwx.Mentions `json:"mentions"`
}

type WecomMarkdownMessage struct {
	MsgType  string            `json:"msgtype"`
	MarkDown WecomMarkdownBody `json:"markdown"`
}

type WecomMarkdownBody struct {
	Content string `json:"content"`
}

type Time time.Time

const (
	timeFormat = "2006-01-02 15:04:05"
)

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	if string(data) == "null" {
		return nil
	}
	now, err := time.ParseInLocation(`"`+timeFormat+`"`, string(data), time.Local)
	*t = Time(now)
	return
}

func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(timeFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, timeFormat)
	b = append(b, '"')
	return b, nil
}

func (t Time) String() string {
	return time.Time(t).Format(timeFormat)
}
