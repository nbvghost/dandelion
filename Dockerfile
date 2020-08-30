FROM alpine

WORKDIR /app

#COPY cert cert
COPY view view
COPY resources resources
COPY cert cert
COPY gweb.json gweb.json
COPY data.json data.json
COPY dandelion dandelion

#RUN ln -s ./niyuan-server /usr/bin/niyuan-server

EXPOSE 8080
EXPOSE 443

CMD ["./dandelion"]