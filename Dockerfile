# GoReleaser v2 multi-arch Dockerfile
# Binary is pre-built by GoReleaser and placed in platform-specific directories
FROM alpine:3.23

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata sqlite \
    && addgroup -g 1000 notebook \
    && adduser -D -u 1000 -G notebook notebook

# Copy pre-built binary from GoReleaser (platform-specific path)
ARG TARGETPLATFORM
COPY ${TARGETPLATFORM}/notebook /usr/local/bin/notebook

# Create data directory for SQLite database
RUN mkdir -p /data && chown notebook:notebook /data

# Use non-root user for security
USER notebook

# Volume for persistent data
VOLUME /data

# Working directory
WORKDIR /data

# Expose default port
EXPOSE 8080

ENTRYPOINT ["notebook"]
CMD ["--dev-listen", ":8080", "--db", "/data/notebook.db"]
