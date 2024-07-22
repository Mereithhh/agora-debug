package main

import (
	"encoding/json"
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

type RecvItem struct {
	DurationMs float64 `json:"durationMs"`
	CostMs     int     `json:"costMs"`
}
type RecvData struct {
	Data []RecvItem
}

func (r *RecvData) Pop() RecvItem {

	item := r.Data[0]
	r.Data = r.Data[1:]
	return item
}

func loadRecvData() RecvData {
	var data []RecvItem
	// 从文件中加载数据
	bs, _ := os.ReadFile("split.json")
	_ = json.Unmarshal(bs, &data)
	return RecvData{
		Data: data,
	}
}

func SimlateRecvTTSData() (ttsChan chan TtsFrame) {
	ttsChan = make(chan TtsFrame, 10)
	go func() {
		// 完全模拟接收到的数据到达时间
		pcmData := LoadPcmData()
		recvData := loadRecvData()
		for {
			if len(recvData.Data) == 0 {
				close(ttsChan)
				break
			}
			info := recvData.Pop()
			time.Sleep(time.Duration(info.CostMs) * time.Millisecond)
			recvDataInBytes := int(info.DurationMs * 32)
			if len(pcmData) < recvDataInBytes {
				close(ttsChan)
				break
			}
			thisData := pcmData[:recvDataInBytes]
			pcmData = pcmData[recvDataInBytes:]
			ttsChan <- NewTTSFrame(thisData)
		}
	}()
	return
}
