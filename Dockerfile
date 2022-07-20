FROM gcr.io/distroless/base-debian11
COPY syoi-access /
ENTRYPOINT [ "syoi-access" ]
