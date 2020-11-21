### Build
# Backend
FROM golang:1.15-alpine3.12 AS backend-builder

RUN mkdir -p /happiness-door-slack-bot
WORKDIR /happiness-door-slack-bot

COPY . .

RUN go build -o out/happiness-door-slack-bot

# Frontend
FROM node:15.2.1-alpine3.12 AS frontend-builder

RUN mkdir -p /happiness-door-slack-bot-frontend
WORKDIR /happiness-door-slack-bot-frontend

COPY ./frontend .

RUN npm install
RUN npm run build

### Runtime
FROM alpine:3.12

WORKDIR /app

COPY --from=backend-builder /happiness-door-slack-bot/out/happiness-door-slack-bot ./happiness-door-slack-bot
COPY --from=frontend-builder /happiness-door-slack-bot-frontend/dist ./happiness-door-slack-bot-frontend

ENV WEB_FOLDER=./happiness-door-slack-bot-frontend

EXPOSE 8080
CMD ./happiness-door-slack-bot