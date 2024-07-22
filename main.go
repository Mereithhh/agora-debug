package main

import (
	"fmt"
	"os"
	"time"

	agora "github.com/Mereithhh/agora-go-sdk"
)

const (
	appId   = "f05303efca1848b2851bfd609d6d7fd4"
	token   = "007eJxTYGj8oVy14qKdpFSiX01P+/mCV0e/drQ+kjZ2/vr111XPlA8KDGkGpsYGxqlpyYmGFiYWSUYWpoZJaSlmBpYpZinmaSkm++/MTWsIZGQ423aXkZEBAkF8FgZDQ0MDBgYAIAoiaw=="
	channel = "1110"
	uid     = "110"
)

func init() {
	removeFiles()
	InitAgora()
}

func removeFiles() {
	_ = os.Remove("agora.log")
	_ = os.Remove("agoraapi.log")
	_ = os.Remove("agoradns.dat")
	_ = os.Remove("agora-debug")
	_ = os.Remove("crash_context_v1")
	_ = os.Remove("xdump_confg")
}

func InitAgora() {
	// 初始化 Agora SDK
	svcCfg := &agora.AgoraServiceConfig{
		AppId:         appId,
		AudioScenario: agora.AUDIO_SCENARIO_CHORUS,
		LogPath:       "agora.log",
		LogSize:       1024 * 1024 * 10,
	}
	agora.Init(svcCfg)
}

func NewConnection() (*agora.RtcConnection, error) {
	waitUserJoined := make(chan struct{})
	conCfg := agora.RtcConnectionConfig{
		SubAudio:       true,
		SubVideo:       false,
		ClientRole:     1,
		ChannelProfile: 1,

		SubAudioConfig: &agora.SubscribeAudioConfig{
			SampleRate: 16000,
			Channels:   1,
		},
	}
	conHandler := agora.RtcConnectionEventHandler{
		OnReconnecting: func(con *agora.RtcConnection, info *agora.RtcConnectionInfo, reason int) {
		},
		OnConnected: func(con *agora.RtcConnection, info *agora.RtcConnectionInfo, reason int) {
			// do something
			fmt.Println("Connected")
		},
		OnUserJoined: func(con *agora.RtcConnection, uid string) {
			fmt.Println("User joined: ", uid)
			waitUserJoined <- struct{}{}
		},
	}
	conCfg.ConnectionHandler = &conHandler
	conCfg.AudioFrameObserver = &agora.RtcConnectionAudioFrameObserver{
		OnPlaybackAudioFrameBeforeMixing: func(con *agora.RtcConnection, channelId string, userId string, frame *agora.PcmAudioFrame) {
			// fmt.Println("OnPlaybackAudioFrameBeforeMixing")
		},
	}
	con := agora.NewConnection(&conCfg)
	nearindump := "{\"che.audio.frame_dump\":{\"location\":\"all\",\"action\":\"start\",\"max_size_bytes\":\"120000000\",\"uuid\":\"123456789\",\"duration\":\"1200000\"}}"
	setResult := con.SetParameters(nearindump)
	if setResult != 0 {
		return nil, fmt.Errorf("SetParameters failed: %d", setResult)
	}
	result := con.Connect(token, channel, uid)
	if result != 0 {
		return nil, fmt.Errorf("Connect failed: %d", result)
	}
	<-waitUserJoined
	return con, nil
}

func main() {

	conn, err := NewConnection()
	if err != nil {
		fmt.Println(err)
		return
	}
	senderContext := NewPcmSender(conn)

	ttsDataChannel := SimlateRecvTTSData()
	for {
		ttsFrame, ok := <-ttsDataChannel
		if !ok {
			fmt.Println("TTS data channel closed, all data received.")
			break
		}
		senderContext.SendPcmData(ttsFrame.data)
	}
	time.Sleep(time.Second * 151)
	senderContext.Sender.Stop()
	conn.Disconnect()
	conn.Release()
}
