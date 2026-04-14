FROM node:20-alpine AS builder

ARG NPM_REGISTRY=https://registry.npmjs.org

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci --registry $NPM_REGISTRY

COPY . .
RUN npm run build

FROM node:20-alpine

ARG NPM_REGISTRY=https://registry.npmjs.org

WORKDIR /app

COPY --from=builder /app/build ./build
COPY --from=builder /app/package.json /app/package-lock.json ./
RUN npm ci --omit=dev --registry $NPM_REGISTRY

EXPOSE 3000

ENV PORT=3000
ENV HOST=0.0.0.0

CMD ["node", "build"]
