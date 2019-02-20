# TokenHolders

1. Migration
open /configs/.default
setup required params:

deployed token addres
TOKEN_ADDRESS=0x57aD67aCf9bF015E4820Fbd66EA1A21BED8852eC

don't change it
RPC_PORT=https://mainnet.infura.io/jV2i3C9g4hfww8EoSpHs

airdrop private key
PRIVATE_KEY=

listen from block
FROM_BLOCK=4941081
listen to block
LAST_BLOCK=7235466

Database must be postgres. Write there connection params.

#DATABASE
DATABASE_USER=postgres
DATABASE_PASSWORD=qwertyui
DATABASE_NAME=tk_db
DATABASE_HOST=0.0.0.0
DATABASE_PORT=5432

Now navigate to /cmd/migration
Write:
go run main.go

Now your DB contains required table.

2. Listener
Navigate to /cmd
Write:
go run main.go

Wait about 5 hours, because there's a lot of transactions. 
Then you'll see "FinalCheck action completed".

3. Distribution
Find in table prewious owner and ecosystem address.
Remove them from table.
Navigate to /cmd/distribution
Write:
go run main.go
Then wait about 10 hours again
Then you'll see "Airdrop completed". 
