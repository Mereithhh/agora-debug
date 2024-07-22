# 拖音问题复现

## 问题描述
在播放音频的过程中，会出现咔的异响或结巴。

可以确定，业务方的tts合成速度是一定跟得上的。

在这个构建的 case 中，模拟服务端合成 tts 的规则是，每 300ms 输出一段 400ms 的音频，但仔细听的话还是能偶尔听到异响/结巴。

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
<audio controls>
  <source src="case/1.wav" type="audio/wave">
Your browser does not support the audio element.
</audio>

场外录音，见 `case/1.wav`，在第 8 秒左右，可以听到咔的一声。

声网日志，见 `case/1agora.log`，没发现 no enough 的日志。

发送的音频就是 `audio.wav`

