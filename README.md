# monodotenv
Generate .env files based on single source configuration file in YAML

# run
./monodotenv use local -f

# features
extends workspaces
fail when no workspace value found
warning and fail if there is existing .env
force option to overrite .env
don't create folders, only .env
add suffix: .env.production
add encrypted gcm
user file

options
value ok
user file ok
encrypted ok
extends ok

sanitize suffix and paths