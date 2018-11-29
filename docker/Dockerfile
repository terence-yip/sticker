FROM debian:stretch-20181112-slim
RUN apt-get update
RUN apt-get install -y \
  sass \
  golang \
  git

ENV GOPATH /app/gopath
RUN mkdir -p $GOPATH
RUN go get -d -u gobot.io/x/gobot/...
 
#RUN cd /app/sticker/scripts && ./generateCss.sh
#RUN cd /app/sticker/scripts && ./buildApp.sh
#RUN cd /app/sticker/staging && ./app
