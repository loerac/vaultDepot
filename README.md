# vaultDepot
Keep your passwords, API keys, etc in a vault. It's like Lastpass, Dashlane, etc but only on your local database.

Your passwords are encrypted in the database and can only be decrypted once you login into your account.

**NOTE:** Still in progress... frontend work is not my forte and is taking me a bit of time.

## Installation
Install some dependencies, `go get -u github.com/loerac/vaultDepot`

Create a database to store your passwords, I have only tested with Postgres so I am not sure how well others will work out.

## Run
`go run main.go`

After it connects to your database, navigate to http://localhost:3000 to create your account and items for your vault.

## Access database
To access your database from where ever you go, you can set up port forwading on modem or use [Dataplicity](https://www.dataplicity.com/) on your machine.
