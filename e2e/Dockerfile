FROM node:20-alpine

ARG NPM_REGISTRY=https://registry.npmjs.org

RUN apk add --no-cache chromium nss freetype harfbuzz ca-certificates ttf-freefont

ENV PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD=1
ENV PLAYWRIGHT_CHROMIUM_EXECUTABLE_PATH=/usr/bin/chromium-browser

WORKDIR /e2e

COPY package.json package-lock.json ./
RUN npm ci --registry $NPM_REGISTRY

COPY . .
