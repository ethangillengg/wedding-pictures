run: build
	@./bin/app

build:
	@go build -o bin/app .


# run templ generation in watch mode to detect all .templ files and 
# re-create _templ.txt files on change, then send reload event to browser. 
# Default url: http://localhost:7331
templ:
	@go run github.com/a-h/templ/cmd/templ@latest generate --watch --proxy="http://localhost:3000" --open-browser=false -v

# run air to detect any go file changes to re-build and re-run the server.
server:
	@go run github.com/cosmtrek/air@v1.51.0 \
	--build.cmd "go build --tags dev -o tmp/bin/main ." --build.bin "tmp/bin/main" --build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

# watch for any js or css change in the assets/ folder, then reload the browser via templ proxy.
sync_assets:
	go run github.com/cosmtrek/air@v1.51.0 \
	--build.cmd "go run github.com/a-h/templ/cmd/templ@latest generate --notify-proxy" \
	--build.bin "true" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_dir "public" \
	--build.include_ext "js,css"

tailwind:
	tailwindcss -i views/css/app.css -o ./public/styles.css --watch

# start the application in development
dev:
	@make -j4 templ server tailwind sync_assets


# build the application for production. This will compile your app
# to a single binary with all its assets embedded.
build:
	@tailwindcss -i views/css/app.css -o ./public/styles.css
	@go build -o bin/app_prod cmd/app/main.go
	@echo "compiled you application with all its assets to a single binary => bin/app_prod"
