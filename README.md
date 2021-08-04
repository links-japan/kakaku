Kakaku oracle service


How to deploy:

1. Build docker image

```
docker build . -t kakaku:latest
```

2. Config setup

Copy config.exmaple.yaml to config.yaml
```
mkdir config
cp config.exmaple.yaml config/config.yaml
```
Fill in all the necessary fields in config.yaml, maily the `mixin` and `db` section.

For `mixin` section you need to register a bot at Mixin developer site: https://developers.mixin.one/dashboard
For `db` section  you need to setup a database instance, here we recommend MySQL.


3. Insert intial asset pairs to database

SQL:
```
insert into assets (base, quote, source, price, term, type) values ("BTC", "JPY", "Coinbase", 0, 1, "Variable"), ("ETH", "JPY", "Coinbase", 0, 1, "Variable"), ("JPYC", "JPY", "", 1, 1, "Const");
```

4. Start service

```
docker run -d --name kakaku -v "$(pwd)"/config:/kakaku/config --env KAKAKU_CONFIG_PATH=/kakaku/config -p 8080:50051/tcp kakaku:latest
```

Now the kakaku service will be running at port :8080
