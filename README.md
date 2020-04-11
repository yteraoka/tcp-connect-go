# Tool to reproduce SYN dropping

- A reason for unexplained connection timeouts on Kubernetes/Docker  
https://tech.xing.com/a-reason-for-unexplained-connection-timeouts-on-kubernetes-docker-abd041cf7e02

- External Source Network Address Translation (SNAT)  
https://docs.aws.amazon.com/eks/latest/userguide/external-snat.html

```
Usage: ./tcp-connect-go [-timeout N] [-count M] [-parallel X] [-threshold Y] [-verbose] server:port
```

```
❯ ./tcp-connect-go -h
Usage of ./tcp-connect-go:
  -count int
    	Number of connect (default 1)
  -parallel int
    	Number of go routines (default 1)
  -threshold int
    	Duration time threshold in millisecond. Show duration time if over this (default 200)
  -timeout int
    	Connect timeout in second (default 10)
  -verbose
    	Enable verbose mode (show duration time every time)
```
