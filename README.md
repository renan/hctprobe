# hctprobe

Probes a HTTP endpoint and creates a TCP Server to indicate a healthy state.

Usage:
```shell
docker pull renang/hctprobe:latest
docker run --publish 8080:8080 renang/hctprobe:latest http://numbersapi.com/42

nc localhost 8080
```
