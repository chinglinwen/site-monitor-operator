FROM harbor.haodai.net/base/alpine:3.7cgo
WORKDIR /app

MAINTAINER wenzhenglin(http://g.haodai.net/wenzhenglin/site-monitor-operator.git)

COPY cmd/manager/site-monitor-operator /app

CMD /app/site-monitor-operator

ENTRYPOINT ["./site-monitor-operator"]

# EXPOSE 8080