FROM localhost:5000/gonode:latest AS build-app
USER root
ENV GO_ENV=production
ENV NODE_ENV=production
COPY . /home/
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
    cd /home/; \
    # build frontend
    npm install -g @angular/cli @nestjs/cli; \
    npm install @types/jest @types/mocha; \
    npm install --save-dev @nestjs/testing; \
    npm install --verbose; \
    [[ -d /export/client/ ]] || { \
        mkdir -pv /export/client/; \
        npm run client:build; \
    }; \
    # build backend
    [[ -d /export/backend/ ]] || { \
        mkdir -pv /export/backend/; \
        npm run backend:build; \
    }; \
    go mod tidy; \
    go get -v ./...; \
    go install -v ./...; \
    go mod vendor; \
    [[ -f /export/api ]] || { \
        go build -ldflags="-s -w" -mod=vendor -o /export/api ./cmd/api/main.go; \
    }; \
    [[ -f /export/downloader ]] || { \
        go build -ldflags="-s -w" -mod=vendor -o /export/downloader ./cmd/downloader/main.go; \
    }; \
    [[ -f /export/encoder ]] || { \
        go build -ldflags="-s -w" -mod=vendor -o /export/encoder ./cmd/encoder/main.go; \
    }; \
    [[ -f /export/file ]] || { \
        go build -ldflags="-s -w" -mod=vendor -o /export/file ./cmd/file/main.go; \
    };

FROM scratch
COPY --from=build-app /export/ /
ENTRYPOINT ["/export/"]
