FROM centos:8

MAINTAINER Fabian Salamanca <fabs@koalatechie.com>

ENV WDIR /opt/fabsgoblog

RUN mkdir $WDIR
RUN echo $WDIR

COPY views /opt/fabsgoblog/views
COPY build/server /opt/fabsgoblog/
COPY .env /opt/fabsgoblog/
COPY start.sh /opt/fabsgoblog/

RUN chmod 755 /opt/fabsgoblog/server

RUN ls -lrth $WDIR

EXPOSE 8019

CMD ["/bin/bash","/opt/fabsgoblog/start.sh"]

