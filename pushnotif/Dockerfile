#BUILDER
FROM golang as builder 
RUN mkdir /build 
ADD . /build 
WORKDIR /build
RUN go build -o app


#ACTUAL
FROM foodora/debian

#copy 
COPY --from=builder /build/app /usr/local/bin/app
RUN chmod +x /usr/local/bin/app
RUN mkdir /app
ADD static /app
WORKDIR /app
EXPOSE 8090
CMD /usr/local/bin/app