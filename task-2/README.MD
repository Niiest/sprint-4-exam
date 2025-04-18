# Задание 2

Переходим в bash-среду брокера kafka-1
- `docker exec -it kafka-1 bash`
- `cd /bin/`

## Тестовые сценарии

Выдаём права консьюмера на группу
- `kafka-acls --authorizer-properties zookeeper.connect=zookeeper:2181 --add --allow-principal User:consumer --operation Read --operation Describe --group my-consumer-group`

### Topic-1

Создаём топик 1
- `kafka-topics --bootstrap-server localhost:9091 --topic topic-1 --create --partitions 3 --replication-factor 2`
- `kafka-topics --bootstrap-server localhost:9091 --topic topic-1 --describe`

Выдаём права для продюсера и консьюмера
- `kafka-acls --authorizer-properties zookeeper.connect=zookeeper:2181 --add --allow-principal User:producer --operation Describe --operation Write --topic topic-1`
- `kafka-acls --authorizer-properties zookeeper.connect=zookeeper:2181 --add --allow-principal User:consumer --operation Describe --operation Read --topic topic-1`
- `kafka-acls --authorizer-properties zookeeper.connect=zookeeper:2181 --add --allow-principal User:consumer --operation Read --operation Describe --group my-consumer-group`

В соседнем терминале запускаем продюсера
- `cd ./task-2/producer`
- `go run .`

В соседнем терминале запускаем консьюмера
- `cd ./task-2/consumer`
- `go run .`

Видим в stdout консьюмера, что данные получены

### Topic-2

Создаём топик 2
- `kafka-topics --bootstrap-server localhost:9091 --topic topic-2 --create --partitions 3 --replication-factor 2`
- `kafka-topics --bootstrap-server localhost:9091 --topic topic-2 --describe`

Выдаём права для продюсера
- `kafka-acls --authorizer-properties zookeeper.connect=zookeeper:2181 --add --allow-principal User:producer --operation Describe --operation Write --topic topic-2`

Выдаём права для консьюмера без права на чтение
- `kafka-acls --authorizer-properties zookeeper.connect=zookeeper:2181 --add --allow-principal User:consumer --operation Describe --topic topic-2`

В соседнем терминале запускаем продюсера
- В `main.go` продюсера меняем переменную topic на значение `topic-2`
- `go run .`

В соседнем терминале запускаем консьюмера
- В `main.go` консьюмера меняем переменную topic на значение `topic-2`
- `go run .`

В stdout консьюмера данные не появляются

Выдаём право на чтение для консьюмера
- `kafka-acls --authorizer-properties zookeeper.connect=zookeeper:2181 --add --allow-principal User:consumer --operation Read --topic topic-2`

В stdout консьюмера данные появились
