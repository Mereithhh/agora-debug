# 拖音问题复现

## 问题描述
在播放音频的过程中，会出现咔的异响或结巴。

可以确定，业务方的tts合成速度是一定跟得上的。

在这个构建的 case 中，模拟服务端合成 tts 的规则完全模拟了服务端的到达时许，具体的时许看 `split.log` 和 `split.json`，模拟逻辑见 `sim_recv_pcm.go`。


## 复现问题

打开 webdemo([https://webdemo.agora.io/example/basic/basicVoiceCall/index.html](https://webdemo.agora.io/example/basic/basicVoiceCall/index.html))，连接频道。

> 频道：  1110   
> uid:   123    
> token(临时的): 007eJxTYGj8oVy14qKdpFSiX01P+/mCV0e/drQ+kjZ2/vr111XPlA8KDGkGpsYGxqlpyYmGFiYWSUYWpoZJaSlmBpYpZinmaSkm++/MTWsIZGQ423aXkZEBAkF8FgZDQ0MDBgYAIAoiaw==

运行项目后，在 webdemo 接收来自 110 的音频，可以听到结巴的现象。
```shell
go run .
```

## 一个 case
场外录音，见 `case/1.m4a`，在第 6 秒左右，第 35 秒左右，都非常奇怪。(最后1s的异常是我强制退出服务端导致的，忽略就行)

声网日志，见 `case/1agora.log`，没发现 no enough 的日志。

发送的音频就是 `audio.wav`

