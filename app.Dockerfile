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
        upx \
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
