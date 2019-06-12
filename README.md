Simple List Creation Bot for Telegram.

PreReqs:

 Install GO 1.12

 Install bitcask cmdline tool
 //go get -u -v github.com/prologic/bitcask
 //cd $GOPATH/src/github.com/prologic/bitcask/cmd/bitcask
 //go install
 

Setup:

Set the TELEGRAMTOKEN environment variable.
Run the bot.
Send "/admin" to the bot.
Record your user id from the console and kill the bot.
//bitcask -p db/admins set $YOURUSERID 1
Restart Bot

Warning:
 No groups or users are whitelisted by default use the /admin commands to add them.

Usage:

/list $LISTNAME -- Prints list items
/list $LISTNAME $ENTRY -- Adds an entry
/admin add group -- Adds current group to whitelist
/admin add admin $ADMIN -- $ADMIN should be an @USER entity, Plaintext will not work. Adds users to admin group (UNTESTED)
