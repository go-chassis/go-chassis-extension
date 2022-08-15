# What is protocol plugin
protocol plugin collocate your implementation to go chassis runtime
a plugin could be any protocol implementation: grpc, gin, fiber etc or even custom protocol.

# Why Collocate your fiber implementation?
with go chassis runtime your native impl can get benefits:
- graceful stop
- [service discovery](../registry): register your server to discovery service such as eureka, servicecomb
- [multiple server in one microservice](https://go-chassis.readthedocs.io/en/latest/user-guides/protocols.html)