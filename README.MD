# tgRandoBot is a telegram bot that generates links to a random channel post

It is pretty self-explanatory. How does it work? You need to forward a message  
from the channel. Then bot will return the link to channel's random message.  
If you just write a basic message to the bot, bot will send default message about  
wanting a forward message instead.

# How to run

1. Clone the project to folder where you will store the bot.
```
https://github.com/rogudator/tgRandoBot.git
```
2. Open the folder.
```
cd tgRandoBot
```
3. Create a .env file containing the token from @botfather . The contents of the .env should be:
```
ID_TOKEN=yourtoken
```
4. Run the bot
```
go run main.go
```