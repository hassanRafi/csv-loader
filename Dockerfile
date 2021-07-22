
FROM golang:1.15.14
WORKDIR /
COPY csv-data /csv-data
COPY /main /main
RUN chmod +x /main
EXPOSE 8000
ENTRYPOINT [ "/main" ]