FROM localhost:5000/gonode:latest AS build-app
USER root
ENV GO_ENV=production
ENV NODE_ENV=production
COPY . /export/
RUN apk update && apk add --no-cache \
        alpine-sdk \
        coreutils \
        wget \
        openssh \
        git \
        build-base \
        apk-tools \
        alpine-conf; \
    \
    cd /export/; \
    # build frontend
    npm install -g @angular/cli @nestjs/cli @nestjs/testing; \
    npm install @types/jest @types/mocha; \
    npm install --verbose; \
    [ -d /export/app/client/ ] || { \
        mkdir -pv /export/app/client/; \
        npm run client:build; \
    }; \
    # build backend
    [ -d /export/app/backend/ ] || { \
        mkdir -pv /export/app/backend/; \
        npm run backend:build; \
    }; \
    go mod tidy; \
    go get -v ./...; \
    go install -v ./...; \
    go mod vendor; \
    [ -f /export/app/api ] || { \
        go build -ldflags="-s -w" -mod=vendor -o /export/app/api ./cmd/api/main.go; \
    }; \
    [ -f /export/app/downloader ] || { \
        go build -ldflags="-s -w" -mod=vendor -o /export/app/downloader ./cmd/downloader/main.go; \
    }; \
    [ -f /export/app/encoder ] || { \
        go build -ldflags="-s -w" -mod=vendor -o /export/app/encoder ./cmd/encoder/main.go; \
    }; \
    [ -f /export/app/file ] || { \
        go build -ldflags="-s -w" -mod=vendor -o /export/app/file ./cmd/file/main.go; \
    };

FROM scratch
COPY --from=build-app /export/app/ /
ENTRYPOINT ["/export/app/"]
