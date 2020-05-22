# happiness-door-slack-bot
Simple Slack bot for measuring https://management30.com/practice/happiness-door/

## Deployment
You need to set up a few ENV variables for the app to work
- `DATABASE_URL` - DB url (e.g. postgres://user:passwd@host:port/dbname)
- `SLACK_TOKEN` - OAuth Slack token which can be found on the app configuration page in section `Install App` 

## Setting up slack application
#### Required permissions for the Slack app are
- `chat:write` - To paste and replace message in chat
- `commands` - To register the /happiness-door command
- `users.profile:read` - To read user profile pictures and using them in voting

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