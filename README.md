# observer.mos.ru data&tools 

<img width="1399" alt="Снимок экрана 2021-09-23 в 07 33 01" src="https://user-images.githubusercontent.com/15847176/134454923-1c38248d-670a-4268-9108-edef00293464.png">


## About


Небольшая CLI для самостоятельного пересчета голосов ДЭГ по выгрузкам сайта https://observer.mos.ru/all/. 
В ней нет реализации проверки целостости самого блокчейна и цифровой подписи голосов. Только их расшифровка.

### Примеры Расшифровки

Пример расшифровки голосования:

- [по одномандатному избирательному округу](data/single-mandate/result.zip)
- [по федеральному избирательному округу](data/consignment/result.zip)

**Графики по данным в сервисе Datalens**

Доступны по ссылке: https://datalens.yandex/pzbi5cq6nd08i

### Как это работает

На сайте https://observer.mos.ru/all/ доступна выгрузка (pg_dump) всех транзакций голосования. 
Ее можно скачать и загрузить через обыкновенный `psql` в свою базу PosgreSQL (потребуется около 20 GB места).

В выгрузке есть две основные таблицы:

- `decrypted_ballots` - расшифрованные самой системой ДЭГ голоса за кандидатов
- `transactions` - исходный лог транзакций

Транзакции делятся на следующие группы по полю `method_id`:

- `10` - Завершение голосования с результатом
- `9` - Расшифровка бюллетеня
- `8` - **Публикация ключа расшифровки**
- `7` - Остановка приёма бюллетеней
- `6` - **Приём бюллетеня (голос)**
- `5` - Проверка доступа голосующего
- `4` - Выдача бюллетеня код метода
- `2` - Системная транзакция
- `1` - Регистрация избирателей
- `0` - **Регистрация кандидатов**

Для "домашнего" пересчета достаточно рассмотреть транзакции - `0, 6, 8`.

### Перечень кандидатов 
В нулевой транзакции хранится перечень кандидатов (их идентификаторы) - имена и номер округа.

```sql
SELECT hash, payload FROM transactions WHERE method_id = 0
```

```json
  "ballots_config": [
{
"district_id": 196,
"options": {
"111906259": "БАЖЕНОВ Тимофей Тимофеевич",
"122700157": "ВИХАРЕВА Эльвира Владимировна",
"133291000": "ИВАНОВА Елена Ивановна",
"141308253": "КРЮКОВ Алексей Сергеевич"
}
}
]
```

### Приём бюллетеня

В транзациях `6` хранится зашифрованный выбор избирателя `encrypted_message` (идентификатор выбранного кандидата).

```sql
SELECT hash, payload FROM transactions WHERE method_id = 6
```

```json
{
  "voting_id": "ea067e1ad71565daff55627e4b35340620d53d644820478ee798e125efe657c2",
  "district_id": 210,
  "encrypted_choice": {
    "encrypted_message": "e7bb71822a92d591ca58532b274726c50bcba5ee22161c3d3d",
    "nonce": "0f1e0f0116a5ab4cf8ba01446edef9718ef6d7bd8f71537a",
    "public_key": "1670df937af8268ce5786e2a2bc4f1080a2f56a1b85727fd34dbf527a7ffab10"
  }
}
```
### Публикация ключа

Мастер ключ для расшифровки публикуется в конце голосования в транзакции `8`.


```sql
SELECT hash, payload FROM transactions WHERE method_id = 8
```

```json
{
  "voting_id": "ea067e1ad71565daff55627e4b35340620d53d644820478ee798e125efe657c2",
  "private_key": "54e3cf70f712b2ff727bde3849772fa811a9d5de796aa7d788d205aa86af04ad",
  "seed": "14901105027823071500"
}
```

Чтобы локально пересчитать голоса достаточно выгрузить все транзакции (`method_id=6`) в csv и запустить утилиту.

```sql
SELECT hash, payload FROM transactions WHERE method_id = 6 > transactions.csv
```

## CLI Installation

```shell
cd voting2021

protoc --go_out=. decryptor/internal/crypto/*.proto 
go build -o ./vote-cli  decryptor/cmd/main.go
```

### CLI Usage

```
Формат  transactions.csv

{hash}, {payload}
{hash}, {payload}
{hash}, {payload}
{hash}, {payload}
```

```shell
./vote-cli {privateKey} transactions.csv
```

```json
./vote-cli 54e3cf70f712b2ff727bde3849772fa811a9d5de796aa7d788d205aa86af04ad examples/votes_2021_truncated.csv

File:  examples/votes_2021_truncated.csv
Key:  54e3cf70f712b2ff727bde3849772fa811a9d5de796aa7d788d205aa86af04ad
Output: out.json
        
2021/09/21 19:38:38 Processed 100 transactions
Candidate ID: 182884641, Votes: 4
Candidate ID: 182247230, Votes: 4
Candidate ID: 149646701, Votes: 1
Candidate ID: 191715167, Votes: 1
Candidate ID: 191070849, Votes: 7
Candidate ID: 193509934, Votes: 4
Candidate ID: 121115662, Votes: 2
Candidate ID: 114042076, Votes: 1
Candidate ID: 221857832, Votes: 1
Candidate ID: 216438542, Votes: 5
Candidate ID: 204289983, Votes: 1
Candidate ID: 162832179, Votes: 5
Candidate ID: 217404809, Votes: 1
Candidate ID: 112984909, Votes: 1
Candidate ID: 153469885, Votes: 1
Candidate ID: 136749451, Votes: 1
Candidate ID: 111906259, Votes: 8
Candidate ID: 115873463, Votes: 11
```

### CLI Output:

```json
[
  {
    "candidate_id": 178404279,
    "votes": 53548
  },
  {
    "candidate_id": 221857832,
    "votes": 6231
  },
  {
    "candidate_id": 182869636,
    "votes": 9522
  },
  {
    "candidate_id": 221327967,
    "votes": 6611
  }
]
```


### Известные проблемы

Не удается расшифровать следующую транзакцию ни моей утилитой, ни публичными инстурментами:
https://observer.mos.ru/all/servers/1/txs

```json
Tx Hash: f07ee512b57bc7d6176592ee6a4ab2526c025af2d57cd9e636c038e61b57db06

Encrypted Message: 3631b09b2513132340eda2d5905756c8dd0d4351b0897e0594
Private Key: 54e3cf70f712b2ff727bde3849772fa811a9d5de796aa7d788d205aa86af04ad
Nonce: 66805830f3a11d610c99b841b22bc71c69227fc6bc9983f1
Public Key: 206c03a67246410a3c992ca13aa005fc68ea70b7ac97f2e4c7e86f2697ba641b
```

### результат расшифровки: 213415260 , КПРФ , Ульянченко

не знаю как написать просто сообщением, поэтому добавил сюда.
P.S. для расшифровки написал свой инструмент на C++

