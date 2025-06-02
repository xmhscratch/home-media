FROM localhost:5000/builder:latest AS build-app
USER root
ENV GO_ENV=production
ENV NODE_ENV=production
COPY . /export/
COPY dist/app/  /export/dist/app/
RUN apk add \
        upx \
        alpine-conf \
    ; \
    cd /export/; \
    \
    if [ ! -d /export/dist/app/frontend/ ] || [ ! -d /export/dist/app/backend/ ]; then \
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
    fi; \
    \
    [ -d /export/dist/app/frontend/ ] || { \
        mkdir -pv /export/dist/app/frontend/; \
        npm run frontend:build &>./dist/app/frontend.out; \
    }; \
    [ -d /export/dist/app/backend/ ] || { \
        mkdir -pv /export/dist/app/backend/; \
        npm run backend:build &>./dist/app/backend.out; \
    }; \
    \
    if [ ! -f /export/dist/app/api ] || [ ! -f /export/dist/app/downloader ] || [ ! -f /export/dist/app/encoder ] || [ ! -f /export/dist/app/file ] || [ ! -f /export/dist/app/tui ]; then \
        go mod tidy; \
        go get -v ./...; \
        go install -v ./...; \
        go mod vendor; \
    fi; \
    for exe in \
		api \
		downloader \
		encoder \
		file \
		tui \
	; do \
		[ -f /export/dist/app/$exe ] || { \
            go build -ldflags="-s -w" -mod=vendor -o /export/dist/app/$exe ./cmd/$exe/main.go &>./dist/app/$exe.out; \
            cp -vrf /export/dist/app/$exe /export/dist/app/$exe.uncompress; \
            upx \
                --no-owner --no-time --no-color -9 -vf --ultra-brute \
                -o/export/dist/app/$exe \
                /export/dist/app/$exe.uncompress \
            ; \
            rm -rf /export/dist/app/$exe.uncompress; \
        }; \
	done;

FROM scratch
COPY --from=build-app /export/dist/app/ /
ENTRYPOINT ["/export/dist/app/"]
