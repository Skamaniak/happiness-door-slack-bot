# happiness-door-slack-bot
Simple Slack bot for measuring https://management30.com/practice/happiness-door/

## Deployment
You need to set up a few ENV variables for the app to work
- `DATABASE_URL` - DB url (e.g. postgres://user:passwd@host:port/dbname)
- `SLACK_TOKEN` - OAuth Slack token which can be found on the app configuration page in section `Install App` 
- `LOG_LEVEL` - Minimal level the app will log. Can be changed runtime. Default Log level is info.
- `WEB_HOST` - The domain of the web frontend part (e.g. <app_name>.herokuapp.com)

## Setting up slack application
#### Required permissions for the Slack app are
- `chat:write` - To paste and replace message in chat
- `chat:write.public` - To post happiness door messages to public channels the bot is not a member of
- `commands` - To register the /happiness-door command
- `users.profile:read` - To read user profile pictures and using them in voting
- `groups:read` and `channels:read` - To find out if the bot can post the message and notify the user about inviting the bot into the channel (typically private ones)

#### Basic app details
- App name: `Happiness Door Bot`
- App icon: can be found in the assets folder

#### Slash command settings
- Command: `/happiness-door`
- Request URL: `https://<host>/rest/v1/happiness-door`
- Short Description: `Create happiness door`
- Usage Hint: `[meeting name]`

#### Interactivity and shortcuts
- Enable interactivity
- Request URL: `https://<host>/rest/v1/happiness-door/interaction`