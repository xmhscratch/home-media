FROM localhost:5000/gonode:latest AS build-app
USER root
ENV GO_ENV=production
ENV NODE_ENV=production
COPY . /export/
RUN apk add --no-cache \
        # alpine-sdk \
        # coreutils \
        # wget \
        # openssh \
        # git \
        # build-base \
        # apk-tools \
        alpine-conf; \
    \
    cd /export/; \
    npm install -g \
        @angular/cli \
        @nestjs/cli; \
    \
    npm install --save-dev \
        @nestjs/testing \
        @nestjs/mapped-types \
        @types/jest; \
    \
    npm install --include=dev; \
    \
    mkdir -pv /export/dist/app/frontend/; \
    npm run frontend:build &>./dist/app/frontend.out; \
    \
    mkdir -pv /export/dist/app/backend/; \
    npm run backend:build &>./dist/app/backend.out; \
    \
    go mod tidy; \
    go get -v ./...; \
    go install -v ./...; \
    go mod vendor; \
    [ -f /export/dist/app/api ] || { \
        go build -ldflags="-s -w" -mod=vendor -o /export/dist/app/api ./cmd/api/main.go &>./dist/app/api.out; \
    }; \
    [ -f /export/dist/app/downloader ] || { \
        go build -ldflags="-s -w" -mod=vendor -o /export/dist/app/downloader ./cmd/downloader/main.go &>./dist/app/downloader.out; \
    }; \
    [ -f /export/dist/app/encoder ] || { \
        go build -ldflags="-s -w" -mod=vendor -o /export/dist/app/encoder ./cmd/encoder/main.go &>./dist/app/encoder.out; \
    }; \
    [ -f /export/dist/app/file ] || { \
        go build -ldflags="-s -w" -mod=vendor -o /export/dist/app/file ./cmd/file/main.go &>./dist/app/file.out; \
    };

FROM scratch
COPY --from=build-app /export/dist/app/ /
ENTRYPOINT ["/export/dist/app/"]
