FROM ubuntu

WORKDIR /srv

COPY ./build /srv/
RUN apt-get update && apt-get install wget -y

RUN chmod +x sai-eth-indexer
CMD ./sai-eth-indexer --debug


EXPOSE 8817
