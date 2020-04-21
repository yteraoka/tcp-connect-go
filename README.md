# Tool to reproduce SYN dropping

- A reason for unexplained connection timeouts on Kubernetes/Docker  
https://tech.xing.com/a-reason-for-unexplained-connection-timeouts-on-kubernetes-docker-abd041cf7e02

- External Source Network Address Translation (SNAT)  
https://docs.aws.amazon.com/eks/latest/userguide/external-snat.html

```
❯ ./tcp-connect-go -h
Usage:
  tcp-connect-go [OPTIONS] [Destination]

Application Options:
  -t, --timeout=        Connect timeout in second. (default: 3)
  -c, --count=          Number of request per threads. (default: 1)
  -p, --parallel=       Number of threads. (default: 1)
      --servername=     Server Name Indication extension in TLS handshake.
  -s, --show-threshold= Duration time threshold in millisecond. Show duration time if over this. (default: 200)
  -v, --verbose         Enable verbose output.
  -V, --version         Show version and exit.
  -S, --sleep=          Sleep time in millisecond after each connect. (default: 0)
  -R, --sleep-random    Randomize sleep time.
      --tls             Do TLS handshake.
  -k, --insecure        Disable certificate verification.

Help Options:
  -h, --help            Show this help message

Arguments:
  Destination:          servername:port
```
