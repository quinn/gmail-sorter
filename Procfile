web: cd "$(pwd -P)" && go tool air -c .air.toml server
tailwind: tailwindcss -i tailwind.css -o internal/web/public/css/styles.css -w=always
templ: cd "$(pwd -P)" && go tool templ generate --watch
tui: cd "$(pwd -P)" && go run main.go tui