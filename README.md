## 🎈log-transfer 



**log-transfer** 项目是 日志采集项目的server端，主要用于日志传输，进行队列缓存，最后进行异步落库操作。



## 🗓特点

+ 基于事件驱动高性能轻量级网络框架`gnet` 的进行开发定制
+ 支持多核多线程
+ 支持`reuseport`套接字IP+端口复用
+ 支持异步读写操作
+ 自定义解码器操作
+ 基于多线程模型的 Event-Loop 事件驱动

## 🎉致谢

感谢 [gnet](https://github.com/panjf2000/gnet)提供的事件驱动的高性能和轻量级网络库

