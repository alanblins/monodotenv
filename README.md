# monodotenv
Generate multiple .env files based on single source configuration file in YAML

# =====WARNING=====
PLEAE DON'T USE IN PRODUCTION. THIS IS STILL UNDER DEVELOPMENT

# Install
- Needs Golang installed
- `go install github.com/alanblins/monodotenv@latest`

# How it works
`monodotenv` will read the file `monodotenv.yaml` and will create .env for a selected environment such as stage, local or testing.

# The problem
Large monorepos present a challenge for engineering teams when it comes to managing environment variables during development. Developers often need to create or update variables, and they must ensure that all team members are aware of these changes to avoid breaking their local environments.

Typically, these environment variables are stored in .env files. However, these files can't be committed to a repository due to security risks, as they often contain sensitive data like passwords or API keys. Additionally, some variables, such as an SDK path, can be user-specific and vary based on a developer's local setup.

This situation forces developers to rely on external documentation to track all the environment variables and their different values for various environments (e.g., local, staging, testing, performance). As a result, developers have to manually manage these variables, leading to significant rework and potential errors whenever a variable is updated.

## How a Solution Like monodotenv can help
A tool like monodotenv aims to solve these issues by storing all environment variables in a single YAML file. This file can be safely committed to a Git repository without exposing sensitive information or user-specific data. This approach streamlines the management of environment variables and ensures that all team members are always synchronized with the latest configurations.

# Use cases

## Basic example
This is just a basic example for a small repository to get familiar with the YAML configuration file.
Create the `monodotenv.yaml` file below.

```yml
environment_variables:
- key: BASE_URL
  environments:
    stage: https://stage.com.myserver
    local: http://localhost:1000
- key: DATABASE_URL
  environments:
    stage: https://stage.database.com
    local: http://localhost:2000
```

The `key` is the environment variable that will be created. It can contain the value from `stage` or `local` that will be sent via CLI. Example:

Run the command below:
```sh
monodotenv use local
```
It will create the a `.env` file with values from environment `local` below:
```
BASE_URL=http://localhost:1000
DATABASE_URL=http://localhost:2000
```

Run the command below:
```sh
monodotenv use stage
```
It will create the .env file with values from environment `stage` below:
```
BASE_URL=https://stage.com.myserver
DATABASE_URL=https://stage.database.com
```

## Creating .env files to different folders
This is useful on monorepo projects with many folders containing their own `.env` files which can contain different or the same environment variables. See the yaml below:
```yaml
environment_variables:
- key: BASE_URL
  environments:
    stage: https://stage.com.myserver
    local: http://localhost:1000
  paths:
  - packages/frontend
- key: AUTH_URL
  environments:
    stage: https://auth.stage.com
    local: http://localhost:3000/auth
  paths:
  - packages/frontend
  - packages/backend
- key: DATABASE_URL
  environments:
    stage: https://stage.database.com
    local: http://localhost:2000
  paths:
  - packages/backend
```
The property `paths` contains a list of folder destination of the `.env` file.
Run:
```sh
monodotenv use local
```
It will create these `.env` files per folder 

*packages/frontend/.env
```
AUTH_URL=http://localhost:3000/auth
BASE_URL=http://localhost:1000
```
* packages/backend/.env
```
AUTH_URL=http://localhost:3000/auth
DATABASE_URL=http://localhost:2000
```
## Add encrypted values
Encrypt a text
```bash
monodotenv enc myapikey1234 mysecretkey0123456789012
```
It will generate output like this:
```bash
Ciphertext: 6abddaad571f4ccc1856bcb5e0ab64728af55bd07646c11af6f8c478
key: mysecretkey0123456789012
nonce: 6a0b01e85109b1c3b6814792
```

Save the `key` and `nonce` into `.monodotenv.secrets.yaml` like below
```yaml
secrets:
 local:
  API_KEY:
  - mysecretkey0123456789012
  - 6a0b01e85109b1c3b6814792
 stage:
  API_KEY:
  - mysecretkey0123456789012
  - 6a0b01e85109b1c3b6814792
```

Add the cipher text into `monodotenv.yaml`.
```yml
environment_variables:
- key: API_KEY
  source: aes-gcm
  environments:
    stage: 6abddaad571f4ccc1856bcb5e0ab64728af55bd07646c11af6f8c478
    local: 6abddaad571f4ccc1856bcb5e0ab64728af55bd07646c11af6f8c478
```
List the environment variables
```
monodotenv use local 
```
The `.env` will be:
```bash
API_KEY=myapikey1234
```

# CLI
## use
Create .env files
```sh
monodotenv use [environment]
```
Ovewrite existing .env files
```sh
monodotenv use [environment]  -f
```
Create .env.testing
```sh
monodotenv use [environment] -s testing
```
## list
List environments
```sh
monodotenv list 
```
List environment variables per environment and per folder destination
```sh
monodotenv list [environment] 
```
## doc
Generate markdown table with enviroment variables
```sh
monodotenv doc 
```
## enc
Encrypt a value
```
monodotenv enc <password> <secret key, 24 or 32 characters>
```
ex:
```
monodotenv enc mypassword c835baf3e8b83e5a
```
Output
```bash
Ciphertext: a9964ffb21cf9c70eff2b179611605e8dafbefac177526958488
key: c835baf3e8b83e5a
nonce: 7a190483b2c1a5226e4f01ab
```
Decrypt a value
```
monodotenv enc -d <ciphertext> <key> <nonce>
```
Ex:
```
monodotenv enc -d a9964ffb21cf9c70eff2b179611605e8dafbefac177526958488 c835baf3e8b83e5a 7a190483b2c1a5226e4f01ab
```
Output
```bash
mypassword
```



# monodotenv.yaml
```
environment_variables: <list>
- key: <environment variable>
  name: <name of the key. Useful for documentation>
  description: <description of the key. Useful for documentation>
  source: <value | user | aes-gcm>
  environments: <key pair. at least one environment required>
    <environment key1>: <environment key2>
    <environment key2>: <environment key2>
  paths: <default: current folder, list>
  - <folder destination 1>
  - <folder destination 2>
```
## environment_variables[].source
### value(default)
Hard coded value from environments
### user
Value from monodotenv.user.file
### aes-gcm
The value in environment will be decrypted using secret available in .monodotenv.secrets.yaml

# features
- Multiple target destinations for .env files. Useful for monorepo projects where needs to repeat the same environment variables on multiple .env files.
- Reuse environment variables with extends environments.
- GCM encrypted values and decrypted while creating .env files
- Generate environment variables per user. Useful for environment variables that is different per user and needs to repeat on multiple .env files