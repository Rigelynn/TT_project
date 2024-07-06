sudo docker rm -f $(sudo docker ps -aq)
sudo docker network prune
sudo docker volume prune
cd fixtures && docker-compose up -d
cd ..
rm tce
go build
netstat -nlp | grep :80 | awk '{print $7}'

docker run \
-p 9002:80 \
--name nginx \
-v /home/nginx/conf/nginx.conf:/etc/nginx/nginx.conf \
-v /home/nginx/conf/conf.d:/etc/nginx/conf.d \
-v /home/nginx/log:/var/log/nginx \
-v /home/nginx/html:/usr/share/nginx/html \
-d nginx:latest

./tce
