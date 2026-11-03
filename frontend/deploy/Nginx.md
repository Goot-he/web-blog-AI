```bash
sudo docker run -d \
  --name nginx \
  --restart unless-stopped \
  --network chenaqi-net \
  -p 80:80 \
  -v /home/docker/frontend/dist:/usr/share/nginx/html:ro \
  -v /home/docker/nginx/nginx.conf:/etc/nginx/conf.d/default.conf:ro \
  docker.m.daocloud.io/library/nginx:1.27-alpine
```