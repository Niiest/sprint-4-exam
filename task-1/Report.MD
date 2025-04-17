# Задание 1

Переходим в bash-среду брокера kafka-0
- `docker exec -it kafka-0 bash`

Создаём топик с 8 партициями и фактором репликации 3
- `kafka-topics.sh --bootstrap-server localhost:9092 --topic balanced_topic --create --partitions 8 --replication-factor 3`

Смотрим описание топика
- `kafka-topics.sh --bootstrap-server localhost:9092 --describe --topic balanced_topic`

Переходим в папку tmp, где можно будет создать `reassignment.json`
- `cd /tmp`

Создаём `reassignment.json`
```
echo '{
    "version": 1,
    "partitions": [
        {"topic": "balanced_topic", "partition": 0, "replicas": [2, 1, 0], "log_dirs": ["any", "any", "any"]},
        {"topic": "balanced_topic", "partition": 1, "replicas": [1, 0, 2], "log_dirs": ["any", "any", "any"]},
        {"topic": "balanced_topic", "partition": 2, "replicas": [1, 2, 0], "log_dirs": ["any", "any", "any"]},
        {"topic": "balanced_topic", "partition": 3, "replicas": [2, 0, 1], "log_dirs": ["any", "any", "any"]},
        {"topic": "balanced_topic", "partition": 4, "replicas": [1, 2, 0], "log_dirs": ["any", "any", "any"]},
        {"topic": "balanced_topic", "partition": 5, "replicas": [2, 0, 1], "log_dirs": ["any", "any", "any"]},
        {"topic": "balanced_topic", "partition": 6, "replicas": [0, 1, 2], "log_dirs": ["any", "any", "any"]},
        {"topic": "balanced_topic", "partition": 7, "replicas": [0, 2, 1], "log_dirs": ["any", "any", "any"]}
    ]
}' > reassignment.json
```

Создаём план для перераспределения партиций
```
kafka-reassign-partitions.sh \
--bootstrap-server localhost:9092 \
--broker-list "1,2,3" \
--topics-to-move-json-file "/tmp/reassignment.json" \
--generate
```

Выполним план для перераспределения партиций
- `kafka-reassign-partitions.sh --bootstrap-server localhost:9092 --reassignment-json-file /tmp/reassignment.json --execute`

Смотрим на результаты выполнения плана перераспределения партиций
- `kafka-topics.sh --bootstrap-server localhost:9092 --describe --topic balanced_topic`
```
Topic: balanced_topic   TopicId: B6-6TRu7TKa_3JopchfqnQ PartitionCount: 8       ReplicationFactor: 3    Configs:
        Topic: balanced_topic   Partition: 0    Leader: 1       Replicas: 2,1,0 Isr: 1,2,0
        Topic: balanced_topic   Partition: 1    Leader: 2       Replicas: 1,0,2 Isr: 2,0,1
        Topic: balanced_topic   Partition: 2    Leader: 0       Replicas: 1,2,0 Isr: 0,1,2
        Topic: balanced_topic   Partition: 3    Leader: 0       Replicas: 2,0,1 Isr: 0,2,1
        Topic: balanced_topic   Partition: 4    Leader: 2       Replicas: 1,2,0 Isr: 2,1,0
        Topic: balanced_topic   Partition: 5    Leader: 1       Replicas: 2,0,1 Isr: 1,0,2
        Topic: balanced_topic   Partition: 6    Leader: 1       Replicas: 0,1,2 Isr: 1,2,0
        Topic: balanced_topic   Partition: 7    Leader: 2       Replicas: 0,2,1 Isr: 2,0,1
```

Во втором терминале останавливаем брокер kafka-1
- `docker stop kafka-1`

Смотрим на результаты остановки относительно топика
- `kafka-topics.sh --bootstrap-server localhost:9092 --describe --topic balanced_topic`
- Для партиций, где лидером был брокер 1, были выбраны новые лидеры автоматичеки
- Брокер 1 перестал быть в списке in-sync replicas (Isr)
```
Topic: balanced_topic   TopicId: B6-6TRu7TKa_3JopchfqnQ PartitionCount: 8       ReplicationFactor: 3    Configs:
        Topic: balanced_topic   Partition: 0    Leader: 2       Replicas: 2,1,0 Isr: 2,0
        Topic: balanced_topic   Partition: 1    Leader: 2       Replicas: 1,0,2 Isr: 2,0
        Topic: balanced_topic   Partition: 2    Leader: 0       Replicas: 1,2,0 Isr: 0,2
        Topic: balanced_topic   Partition: 3    Leader: 0       Replicas: 2,0,1 Isr: 0,2
        Topic: balanced_topic   Partition: 4    Leader: 2       Replicas: 1,2,0 Isr: 2,0
        Topic: balanced_topic   Partition: 5    Leader: 2       Replicas: 2,0,1 Isr: 0,2
        Topic: balanced_topic   Partition: 6    Leader: 0       Replicas: 0,1,2 Isr: 2,0
        Topic: balanced_topic   Partition: 7    Leader: 2       Replicas: 0,2,1 Isr: 2,0
```

Запускаем брокер kafka-1 во втором терминале
- `docker start kafka-1`
- Брокер 1 вернулся в список in-sync replicas, но не стал быть лидером для какой-либо партиции
- Для улучшенной балансировки необходимо ручное переназначение лидеров
```
Topic: balanced_topic   TopicId: B6-6TRu7TKa_3JopchfqnQ PartitionCount: 8       ReplicationFactor: 3    Configs:
        Topic: balanced_topic   Partition: 0    Leader: 2       Replicas: 2,1,0 Isr: 2,0,1
        Topic: balanced_topic   Partition: 1    Leader: 2       Replicas: 1,0,2 Isr: 2,0,1
        Topic: balanced_topic   Partition: 2    Leader: 0       Replicas: 1,2,0 Isr: 0,2,1
        Topic: balanced_topic   Partition: 3    Leader: 2       Replicas: 2,0,1 Isr: 0,2,1
        Topic: balanced_topic   Partition: 4    Leader: 2       Replicas: 1,2,0 Isr: 2,0,1
        Topic: balanced_topic   Partition: 5    Leader: 2       Replicas: 2,0,1 Isr: 0,2,1
        Topic: balanced_topic   Partition: 6    Leader: 0       Replicas: 0,1,2 Isr: 2,0,1
        Topic: balanced_topic   Partition: 7    Leader: 0       Replicas: 0,2,1 Isr: 2,0,1
```
