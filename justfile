run *args:
  go run main.go {{args}}

build:
  go build

dev:
  go tool goreman -b 5001 start
