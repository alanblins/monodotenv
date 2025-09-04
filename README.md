# monodotenv
Generate multiple .env files based on single source configuration file in YAML

# The problem

# Install
- Needs Golang installed
- Clone this repo
- Compile it `go build`
- Create a basic yaml file
- Run it `./monodotenv use local`

# How it works
`monodotenv` will read the file `monodotenv.yaml` and will create .env for different workspace like stage, local, testing, etc.
## Basic example
Create the `monodotenv.yaml` file below.

```yml
environment_variables:
- key: BASE_URL
  workspaces:
    stage: https://stage.com.myserver
    local: http://localhost:1000
- key: DATABASE_URL
  workspaces:
    stage: https://stage.database.com
    local: http://localhost:2000
```

The `key` is the environment variable that will be created. It can contain the value from `stage` or `local` that will be sent via CLI. Example:


Execute the command below:
```sh
./monodotenv use local
```
It will create the .env file with values from workspace `local` below:
```
BASE_URL=http://localhost:1000
DATABASE_URL=http://localhost:2000
```

Execute the command below:
```sh
./monodotenv use stage
```
It will create the .env file with values from workspace `stage` below:
```
BASE_URL=https://stage.com.myserver
DATABASE_URL=https://stage.database.com
```

## Creating .env files to different folders
This is useful on monorepo project where we need .env files on multiple places. See the yaml below:
```yaml
environment_variables:
- key: BASE_URL
  workspaces:
    stage: https://stage.com.myserver
    local: http://localhost:1000
  paths:
  - packages/frontend
- key: AUTH_URL
  workspaces:
    stage: https://auth.stage.com
    local: http://localhost:3000/auth
  paths:
  - packages/frontend
  - packages/backend
- key: DATABASE_URL
  workspaces:
    stage: https://stage.database.com
    local: http://localhost:2000
  paths:
  - packages/backend
```
The property `paths` contains a list of destination of the environment variables.
Execute the command:
```sh
./monodotenv use local
```
It will create different `.env` files at `packages/frontend` and `packages/backend`.
The EV `BASE_URL` will be at `packages/frontend`
The EV `AUTH_URL` will be at `packages/frontend` and `packages/backend`
The EV `DATABASE_URL` will be at `packages/backend`
The final files will be as below:
```
packages/backend/.env
AUTH_URL=http://localhost:3000/auth
DATABASE_URL=http://localhost:2000

packages/frontend/.env
BASE_URL=http://localhost:1000
AUTH_URL=http://localhost:3000/auth
```


# CLI
Create .env files
```sh
./monodotenv use [workspace]
```
Ovewrite existing .env files
```sh
./monodotenv use [workspace]  -f
```
Create .env.testing
```sh
./monodotenv use [workspace] -s testing
```
List workspaces
```sh
./monodotenv list 
```
List environment variables per workspace and per folder destination
```sh
./monodotenv list [workspace] 
```

# monodotenv.yaml
```
environment_variables: <list>
- key: <environment variable>
  source: <default: value, value = hard coded value from workspaces | user = value from monodotenv.user.file | aes-gcm = the value in workspace will be decrypted using secret available in .monodotenv.secrets.yaml>
  workspaces: <key pair. at least one workspace required>
    <workspace key1>: <workspace key2>
    <workspace key2>: <workspace key2>
  paths: <default: current folder, list>
  - <folder destination 1>
  - <folder destination 2>
```

# features
- Multiple target destinations for .env files. Useful for monorepo projects where needs to repeat the same environment variables on multiple .env files.
- Reuse environment variables with extends workspaces.
- GCM encrypted values and decrypted while creating .env files
- Generate environment variables per user. Useful for environment variables that is different per user and needs to repeat on multiple .env files