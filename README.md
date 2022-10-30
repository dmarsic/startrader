# Star Trader

It's a game!

It's not a playable game. I wrote this as an attempt to learn Go.

The game runs a server, and users can interact by sending API requests. See [client.sh](client.sh).

There is no real authentication, but there is a login endpoint, and
a user must login to store the user ID into a cookie. The user must
match one of the users under [data/users](data/users), otherwise there
will be no data associated with them.

## Example

Run the game to start the server.

```bash
go run main.go
```

Create new user.

```bash
curl -X POST --data "name=username" http://localhost:5000/u/new
```

Login to store the user id as a session variable into `cookie.txt`.

```bash
curl -c cookie.txt http://localhost:5000/login?user=username
```

Show my information. We need to reference the cookie file.

```bash
curl -b cookie.txt http://localhost:5000/u/username
```

Note that the root location redirects to showing our user's information.

```bash
curl -Lb cookie.txt http://localhost:5000/
```

Routes are listed in [main.go](main.go).

# License

See LICENSE.

# Author

Dom Marsic, 2022.
