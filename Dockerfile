FROM golang:1.21-alpine AS build

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -o /app/task_manager ./webapp

FROM scratch
COPY --from=build /app/task_manager /app/task_manager

EXPOSE 8080

ENTRYPOINT ["/app/task_manager"]
