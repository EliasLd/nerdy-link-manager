# ----------- Build stage -----------
FROM node:20-alpine AS build

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .
RUN npm run build

# ----------- Runtime stage -----------
FROM docker.io/caddy:2.10-alpine

COPY Caddyfile /etc/caddy/Caddyfile

COPY --from=build /app/build /usr/share/caddy

EXPOSE 80
