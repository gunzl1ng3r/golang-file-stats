FROM alpine:3.4
# installs unrar
RUN apk add --no-cache unrar

# create directory to work with and work in
RUN mkdir -p /files
WORKDIR /files

# execute find and make sure to use ';' instead of \;
ENTRYPOINT [ "/bin/sh","-c","find /files -name '*.rar' -exec unrar e {} ';'" ]
