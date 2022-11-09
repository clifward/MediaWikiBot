WikiBot
---

# Overview

The purpose of this bot is to read frequently the recent action API of MediaWiki-based wiki, and to send them to a Telegram group/channel using a bot.

# Features

The program can run in 2 differents modes :

- job, which has to be triggered frequently by another mecanism, such as kubernetes cronjob or AWS Lambda
- loop, which will manage the job trigger (you should use the docker image for this usecase)

# Deployment

You just have to run `update.sh`
