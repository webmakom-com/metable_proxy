FROM ubuntu

WORKDIR /srv

COPY ./build /srv/

RUN apt-get update && apt-get install -y --reinstall ca-certificates openssl

RUN openssl req -x509 -nodes -newkey rsa:2048 -keyout server.rsa.key -out server.rsa.crt -days 3650 -subj '/CN=localhost'
RUN ln -sf server.rsa.key server.key
RUN ln -sf server.rsa.crt server.crt

RUN chmod +x sai-auth
CMD ./sai-auth start

EXPOSE 8888
EXPOSE 8800
EXPOSE 8803
