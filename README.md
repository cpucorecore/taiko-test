# install foundry
https://book.getfoundry.sh/getting-started/installation

# generate accounts
```bash
cd scripts
bash generate_accounts.sh
```

# transfer 1tether to the account generated above
in the scripts' path:
```bash
bash init_account_balance_with_1tether.sh
```

# do test
```
go run .
```

this will do 20 transfer tx for every account generated above, total: 20txs * 10000 accounts = 20w txs