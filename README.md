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