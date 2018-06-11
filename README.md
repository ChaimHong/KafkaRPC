# kafka-rpc
go rpc over kafka message queue
简易RPC工具，使用kafka作为rpc消息的中间件，强化rpc功能。
让远程服务器挂掉或者负载高的时候，还能继续提供rpc异步回调服务！

# Proto
协议默认是用[gobuf](https://github.com/ChaimHong/gobuf)
也有支持gob的分支

# Cluster
middleware对应的是kafka的consumer group概念，负载高的时候可以启多个进程的

[详看cluster分支](https://github.com/ChaimHong/kfkrpc/tree/cluster)
