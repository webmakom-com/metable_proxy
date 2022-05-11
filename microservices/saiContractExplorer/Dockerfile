FROM ubuntu

WORKDIR /srv

COPY ./build /srv/
RUN apt-get update && apt-get install wget -y

RUN chmod +x sai-contract-explorer
CMD ./sai-contract-explorer

EXPOSE 8817
