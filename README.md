# grpc-kcp-transport

this is a helper function for you can transport grpc request over kcp

"kcp is a lightweight and fast network protocol, it is used for high performance communication between two endpoints."



how to use
```
-> go get -u github.com/Lynnworld/grpc-kcp-transport
```

```go
// server 
address := "localhost:8080"
cfg := &kcp.Config{}
server := grpc.NewServer()
demo.RegisterDemoServer(server, &serverImpl{})
listener, err := kcp.Listener(address, config)
server.Serve(listener)

// client
address := "localhost:8080"
cfg := &kcp.Config{}
cc, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()), kcp.Dialer(cfg))

```
you can view example for grpc server / client demo , it use both TCP/KCP for grpc serving


for more details, please visit:

1. https://github.com/skywind3000/kcp
2. https://github.com/xtaci/kcp-go
3. https://github.com/xtaci/kcptun
