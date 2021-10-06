build-bin:
	echo "::: Building web project -> bin"
	go build -o bin/gh_statistics -v .

run-heroku:
	echo "::: Running web project in heroku local"
	heroku local web

# all-web: build-web run-web
# all-cli: build-cli run-cli