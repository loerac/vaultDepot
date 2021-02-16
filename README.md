# vaultDepot
Keep your passwords, API keys, etc in a vault. It's like Lastpass, Dashlane, etc but only on your local database.

Your passwords are encrypted in the database and can only be decrypted once you login into your account.

**NOTE:** Still in progress... I was making a web portion but frontend work is not my forte and is taking me a bit of time. So for now, only using the command line.

## Installation
Install some dependencies, `go get -u github.com/loerac/vaultDepot`

Create a database to store your passwords, I have only tested with Postgres so I am not sure how well others will work out.

## Run
`go run main.go`

Create an account with a username, password, and secret key. The password is for your account in the database, and the secret key is for encryting and decrypting your passwords.

## Access database
To access your database from where ever you go, you can set up port forwading on modem or use [Dataplicity](https://www.dataplicity.com/) on your machine.
