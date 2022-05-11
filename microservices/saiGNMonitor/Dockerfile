FROM ubuntu

WORKDIR /srv

COPY ./build /srv/
RUN apt-get update && apt-get install wget -y

RUN chmod +x sai-gn-monitor
CMD ./sai-gn-monitor start

EXPOSE 8820
