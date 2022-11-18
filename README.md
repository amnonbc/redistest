# Benchmark of Redis Client Libraries

benchmark pub-sub operations of github.com/go-redis/redis/v8 vs github.com/rueian/rueidis

go-redis baseline
```
% time ./standardlib/standardlib -n 100000
redistest.go:51: got 100000 messages in 5.001932505s

real	0m5.006s
user	0m1.920s
sys	0m3.404s
```

rueian gives similar performance
```
% time ./rueian/rueian -n 100000 -batchsize 1
2022/11/18 17:06:31 got 100000 messages in 5.899845305s

real	0m5.905s
user	0m1.553s
sys	0m4.115s
```

If we batch the PUBLISH requests, we can get orders of magnitude faster

```
% time ./rueian/rueian -n 100000 -batchsize 10
2022/11/18 17:06:35 got 100000 messages in 880.754704ms

real	0m0.886s
user	0m0.362s
sys	0m0.511s
% time ./rueian/rueian -n 100000 -batchsize 1000
2022/11/18 17:06:40 got 100000 messages in 318.437865ms

real	0m0.324s
user	0m0.235s
sys	0m0.020s
```

If we issue the PUBLISH requests concurrently, then we get a 5 x reduction in elapsed time
and a 10x reduction in system time.

```
% time ./rueian/rueian -n 100000 -concurrent=true
2022/11/18 17:07:11 got 100000 messages in 830.2702ms

real	0m0.858s
user	0m0.898s
sys	0m0.337s
```