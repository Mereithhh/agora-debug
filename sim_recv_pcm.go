package main

import (
	"os"
	"time"
)

func LoadPcmData() []byte {
	bs, _ := os.ReadFile("audio.wav")
	bs = bs[44:]
	return bs
}

type TtsFrame struct {
	data []byte
}

func NewTTSFrame(data []byte) TtsFrame {
	return TtsFrame{data: data}
}

func SimlateRecvTTSData() (ttsChan chan TtsFrame) {
	ttsChan = make(chan TtsFrame, 10)
	go func() {
		// 每 300ms 收到一个 400ms 的音频帧
		recvInterval := time.Millisecond * 300
		recvDurationMs := 400
		pcmData := LoadPcmData()
		for {
			if len(pcmData) < recvDurationMs*32 {
				close(ttsChan)
				break
			}
			thisData := pcmData[:recvDurationMs*32]
			pcmData = pcmData[recvDurationMs*32:]
			ttsChan <- NewTTSFrame(thisData)
			time.Sleep(recvInterval)
		}
	}()
	return
}
