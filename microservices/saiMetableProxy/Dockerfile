FROM ubuntu

WORKDIR /srv

COPY ./build /srv/

RUN chmod +x sai-metable-proxy
CMD ./sai-metable-proxy start

EXPOSE 8818
