tailwindcss:
	npx @tailwindcss/cli -i assets/app.css -o public/app.css --watch=always

templ:
	templ generate --watch --proxy="http://localhost:2999" --proxybind="0.0.0.0" --proxyport="3000" --cmd="go run ."

dev:
	-@fuser -k 2999/tcp
	-@fuser -k 3000/tcp
	make -j2 templ tailwindcss