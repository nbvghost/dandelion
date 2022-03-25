FROM alpine

WORKDIR /app

#COPY cert cert


#RUN ln -s ./niyuan-server /usr/bin/niyuan-server

EXPOSE 8080
EXPOSE 443

CMD ["./dandelion"]