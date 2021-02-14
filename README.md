<h1 align="center">twitter-discord-bot (Test Project)</h1>

# About
Mini go script that automatically shares your tweets on Discord channel via webhook.

# Usage
```shell
$ go run main.go -key=<Twitter Key> -secret=<Twitter Secret> -username=<Twitter Username> -displayName=<Webhook Display Name> -avatarUrl=<Webhook Avatar Url> -webhook=<Discord Webhook Url>
```

# Automation
```shell
$ crontab -e
```
then define cron job rule
```text
* * * * * go run main.go -key=<Twitter Key> -secret=<Twitter Secret> -username=<Twitter Username> -displayName=<Webhook Display Name> -avatarUrl=<Webhook Avatar Url> -webhook=<Discord Webhook Url> -dataFolder=/home/bot >> /home/log.txt
```
****
