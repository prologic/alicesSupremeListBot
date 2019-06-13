
## Simple List Creation Bot for Telegram.
### PreReqs:
- Install GO 1.12
- Install bitcask cmdline tool
`go get -u -v github.com/prologic/bitcask`
`cd $GOPATH/src/github.com/prologic/bitcask/cmd/bitcask`
`go install`

### Setup:
1. Set the TELEGRAMTOKEN environment variable.
2. Run the bot.
3. Send `/admin` to the bot.
4. Record your user id from the console and kill the bot.
5. `bitcask -p db/admins set $YOURUSERID 1`
6. Restart Bot

## Warning:
No groups or users are whitelisted by default use the /admin commands to add them.

## Usage:
- `/list $LISTNAME` -- Prints list items
- `/list $LISTNAME $ENTRY` -- Adds an entry
- `/lists` -- Prints public lists. 
- `/admin add group` -- Adds current group to whitelist
- `/admin invite $ADMIN` --  Invites users to admin group.
- `/admin accept` -- Accepts pending admin invite


### Notes:
- Lists that start with `_` are considered private and are not included in `/lists` but they still are accessible if you know the name.
-  `$ADMIN` should be an entity of type mention or text_mention (@USER) when using `/admin add`, Plaintext will not work.
- `/list $LISTNAME $ENTRY` will create the list if it doesn't exist
- `$LISTNAME` is case-insensitive
- Users must accept admin invites due to a telegram bot API limitation, unless the user has no username then they will be immediately added.
## Possible Future Features
- Group Locked Lists
- Delete Lists
- Delete List Entries
- TBD
