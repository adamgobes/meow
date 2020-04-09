release:
	docker-compose build web
	docker tag meow_web registry.heroku.com/statcatmeow/web
	docker push registry.heroku.com/statcatmeow/web
	heroku container:release web --app statcatmeow

dev:
	docker-compose up -d db local