web: cd "$(pwd -P)" && go tool air -c .air.toml server
tailwind: tailwindcss -i tailwind.css -o internal/web/public/css/styles.css -w=always
templ: go tool templ generate --proxy="http://localhost:$PORT" --open-browser=false --watch
