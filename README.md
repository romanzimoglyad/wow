# Words of Wisdom with Proof Of Work

### PoW algorithm

Hashcash https://en.wikipedia.org/wiki/Hashcash

Implemented checks:
+ The recipient's computer calculates the 160-bit SHA-1 hash of the entire string (e.g., "1:20:060408:anni@cypherspace.org::1QTjaYd7niiQA/sc:ePa"). This takes about two microseconds on a 1 GHz machine, far less time than the time it takes for the rest of the e-mail to be received. If the first 20 bits are not all zero, the hash is invalid. (Later versions may require more bits to be zero as machine processing speeds increase.)
+ The recipient's computer checks the date in the header (e.g., "060408", which represents the date 8 Apr 2006). If it is not within two days of the current date, it is invalid. (The two-day window compensates for clock skew and network routing time between different systems.)
+ The recipient's computer checks whether the e-mail address in the hash string matches any of the valid e-mail addresses registered by the recipient, or matches any of the mailing lists to which the recipient is subscribed. If a match is not found, the hash string is invalid.
+ The recipient's computer inserts the hash string into a database. If the string is already in the database (indicating that an attempt is being made to re-use the hash string), it is invalid. 


Restrictions: instead of using a real database I've used simple in-memory storage

### Protocol

Type|Body
#### Type
+ 0 - Quit / Response
+ 1 - GetChallenge
+ 2 - GetMessage
+ 3 - Error

### Example of use


1. Start server:

    `make start-server`

2. Connect to server from terminal:

    `telnet localhost 8801`

3. Send request for challenge:

   ` 1|`

   Response from server can be like:

    `0|{"Ver":1, "Bits":4, "Date":1655118481, "Resource":"127.0.0.1:38784", "Rand":"NzgxMTIxMQ==", "Counter":0}`


4. Count Hash and send request for resource, for example:
`2|{"Ver":1, "Bits":4, "Date":1655118481, "Resource":"127.0.0.1:38784", "Rand":"NzgxMTIxMQ==", "Counter":1213}`
   
   Response from server can be like:

   `0|Word of Wisdom: Do what inspires you. Life is too short not to love the job you do every day.`

5. In case of server error response from server can be like:
   `3|invalid hashcash`
6. Close connection:

    `0|`
## Example of the configuration used

Below is a list of variables that should be used to run the project:

| ENV VAR                            | Example Value | Description                          | Default value |
|------------------------------------|---------------|--------------------------------------|---------------|
| WOW_IP                             | `'0.0.0.0'`   | server host                          | `0.0.0.0`     |
| WOW_PORT                           | `'8801'`      | server port                          | `8801`        |
| WOW_LOGLEVEL                       | `'warn'`      | log level                            | `warn`        |
| WOW_ZERO_NUMBER                    | `'4'`         | number of zeroes in POW              | `4`           |
| WOW_MAX_IT                         | `'1000000'`   | max number of hash count iterations  | `1000000`     |
| WOW_CLIENT_REQUEST_NUMBER          | `'100'`       | request number from client           | `100`         |
| WOW_CLIENT_SEND_INTERVAL_MS        | `'1000'`      | client send interval in milliseconds | `1000`        |

## Make commands

### RUN DEMO (server + client)
make run
### STOP
make stop
### RUN tests
make test
### RUN server
make start-server
### RUN client
make start-client