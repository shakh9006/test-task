# Test Task

# RUN

```
docker-compose up -d --build 
```

OR 
```
make
```

## After compose docker run:
```
docker ps
```

![alt text](docs/1.png)


### To check database rows write:
```
docker exec -it <postgres_hash_id> sh
```
```
psql -U postgres
```
```
\d
```
```
SELECT * FROM numbers;
```

### To check logs write:
```
docker exec -it <numbers_hash_id> sh
```

```
ls -la
cd logs
ls -la

cat info.log
cat error.log
cat debug.log
```

![alt text](docs/2.png)


### Queues:
![alt text](docs/3.png)