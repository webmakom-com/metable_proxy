FROM ubuntu

WORKDIR /srv

COPY ./build /srv/
RUN apt-get update && apt-get install wget -y

RUN chmod +x sai-eth-manager
CMD ./sai-eth-manager start

EXPOSE 8819
