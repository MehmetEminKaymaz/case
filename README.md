# Kurulum

# Gereksinimler :
- Docker
- Goland

# RabbitMQ & PostgreSql : 
- docker run --rm -it -p 15672:15672 -p 5672:5672 rabbitmq:3-management
- docker run --name my-postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=case -e POSTGRES_DB=postgres -p 5432:5432 -d postgres:latest

# Açıklama

Dispatcher, consumer ve api tek projede olmasının sebebi kolaylık olsun diye normalde dağıtık kullanılabilir. Her 2 dakikada bir gerekli kontrol yapılıp data kuyruğa iletilir, consumer kuyruğu okur mesajı ilettiyse asıl kaydı günceller. Basit bir outbox pattern ile 3 app olması gereken yapı tek app olarak bulunuyor. Mesajların statüsü db üzerinden veya aşağıdaki örnek curller ile takip edilebilir.

Kullanım için örnek curller : 

Gönderilecek mesajı oluşturmak için (POST) :

curl --location 'http://localhost:8080/message' \
--header 'Content-Type: application/json' \
--data '{
    "content" : "test",
    "recipient" : "05551111111"
}'

Gönderilmişleri listelemek için (GET): 

curl --location 'http://localhost:8080/messages'

Dispatcher'ı çalıştırmak ve kapatmak için (GET): 

curl --location 'http://localhost:8080/start' && curl --location 'http://localhost:8080/stop'
