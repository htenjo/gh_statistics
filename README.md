# gh_statistics
Repo to authenticate a user and configure basic useful statistics from GitHub

# Current process (not very happy, but it works for now)
1. Run the web project
2. Open localhost:8080
3. If you don't have the privileges, the GitHub OAuth process will start
4. Once the authorization is granted, you can go to the /repos page to persist the repos you want to monitor
5. When you are ready, just click on Send Notification

## How to run
``` bash
$ make build-cli
$ make run-cli sid=<some-gh-profile-id-configured-from-web>
```

## How to install postgres
``` bash 
docker run --name gh_stats_db -e POSTGRES_PASSWORD=admin -p 5432:5432 -d postgres
```

## Heroku commands
| Command | Description |  
|---|---|
|`$ heroku apps:create <APP_NAME>`|Creates a new application in heroku|
|`$ heroku config:set <VAR_NAME>=<VAR_VALUE>>`| Creates a new env var|
|`$ heroku config`| Lists configured env vars|
|`$ heroku addons:create heroku-postgresql:<dbName>` |Create a new PostgresDB|
|`$ heroku pg:psql`|Allows open sql command line|
|`$ heroku logs --tail`| Display current process logs|
|`$ heroku local [web]`| Run the web project in the local machine|