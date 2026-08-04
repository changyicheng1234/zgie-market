curl <http://localhost:8080/debug/pprof/trace?seconds=20>> trace.out

# 在宿主机执行，把 trace.out 从容器里复制出来

docker cp <你的容器名>:/app/trace.out ./trace.out

go tool trace -http=127.0.0.1:8081 trace.out
