FROM debian:bookworm-slim
WORKDIR /app
COPY ./tmp/oe3xorstatus .
COPY ./tmp/pb_public ./pb_public
CMD ["/app/oe3xorstatus", "serve", "--http=0.0.0.0:8090"]