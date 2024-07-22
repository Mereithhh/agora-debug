package main

import (
	"fmt"
	"time"

	agoraservice "github.com/Mereithhh/agora-go-sdk"
)

type PcmSender struct {
	HasSendPcmData bool
	Sender         *agoraservice.PcmSender
	SendAudioCache *AudioCache
	Closed         *BoolValue
}

func NewPcmSender(conn *agoraservice.RtcConnection) *PcmSender {
	sender := conn.NewPcmSender()
	sender.Start()
	return &PcmSender{
		HasSendPcmData: false,
		SendAudioCache: NewAudioCache(),
		Sender:         sender,
		Closed:         NewBoolValue(false),
	}
}

func (l *PcmSender) SendPcmDataAtom(data []byte) error {
	frame := &agoraservice.PcmAudioFrame{
		Data:              data,
		Timestamp:         0,
		SamplesPerChannel: 160,
		BytesPerSample:    2,
		NumberOfChannels:  1,
		SampleRate:        16000,
	}
	result := l.Sender.SendPcmData(frame)
	if result != 0 {
		return fmt.Errorf("SendPcmData failed: %d", result)
	}
	return nil

}

func (l *PcmSender) SendPcmData(data []byte) {
	if !l.HasSendPcmData {
		l.HasSendPcmData = true
		go l.StartTimer(data)
		return
	}
	l.PutPcmDataInTimer(data)
}

func (l *PcmSender) StartTimer(bootstartData []byte) {
	fmt.Printf("[语音通话]loopctx 开始发送音频包, 启动音频大小：%d\n", len(bootstartData))
	l.Sender.ClearSendBuffer()
	l.SendAudioCache.Put(bootstartData)
	sendCount := 0
	var startTime int64 = 0
	onClose := func() {
		l.SendAudioCache.Clear()
	}

	for i := 0; i < 18; i++ {
		err := l.SendPcmDataAtom(l.SendAudioCache.GetSize(320))
		if err != nil {
			fmt.Println("sendPcmDataAtom failed")
		}
		sendCount++
	}

	ticker := time.NewTicker(time.Millisecond * 50)
	defer ticker.Stop()
	startTime = time.Now().UnixMilli()
	onTick := func() {
		shouldSendCount := (time.Now().UnixMilli()-startTime)/10 - int64(sendCount-18)
		hasCount := l.SendAudioCache.Size() / 320
		fmt.Printf("[语音通话]loopctx 应该发送音频包数量:%v， 已发送个数:%v， 剩余个数: %v \n", shouldSendCount, sendCount, hasCount)
		if shouldSendCount > int64(hasCount) {
			fmt.Printf("[语音通话]loopctx 发送音频包数量不够了，还差 %v 个\n", shouldSendCount-int64(hasCount))
		}
		for i := 0; i < int(shouldSendCount); i++ {
			d := l.SendAudioCache.GetSize(320)
			if d == nil {
				i--
				time.Sleep(time.Millisecond * 1)
				continue
			}
			err := l.SendPcmDataAtom(d)
			if err != nil {
				fmt.Printf("[语音通话]loopctx 发送音频包失败，错误信息: %v\n", err)
				break
			}
			sendCount++
		}
	}
	onTick()
	for range ticker.C {
		if l.Closed.Get() {
			fmt.Printf("[语音通话]loopctx 已关闭停止发送音频包\n")
			onClose()
			return
		}
		onTick()
	}
}

func (l *PcmSender) PutPcmDataInTimer(data []byte) {
	l.SendAudioCache.Put(data)
}

func (l *PcmSender) Close() {
	l.Closed.Set(true)
}
