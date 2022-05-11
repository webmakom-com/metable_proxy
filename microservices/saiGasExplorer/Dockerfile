FROM ubuntu

WORKDIR /srv

COPY ./build /srv/
RUN apt-get update && apt-get install wget -y

RUN chmod +x sai-gas-explorer
CMD ./sai-gas-explorer

EXPOSE 8821
