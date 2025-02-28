FROM debian:stable-slim

COPY Chirpy /bin/Chirpy

COPY index.html /index.html

CMD ["/bin/Chirpy"]

ENV PORT=8080
