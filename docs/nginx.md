# nginx

Placeholder for documentation for the nginx reverse proxy.

## Setup

First, setup api on server, then proceed. Server names
will need edited in the nginx conf for your own domain.
You also need an A record for subdomain `api` to port
hosting. This assumes my setup for small, 1 vm instance
rather than a load balancer which should be used for
a production service.

Setup nginx:

```bash
sudo dnf install nginx -y
sudo systemctl enable nginx
sudo systemctl start nginx
sudo systemctl status nginx
```

Copy nginx conf:

```bash
sudo cp api.wellbeingquest.app.conf /etc/nginx/conf.d/
```

Setup firewall ports:

```bash
sudo firewall-cmd --permanent --add-service=http
sudo firewall-cmd --permanent --add-service=https
sudo firewall-cmd --reload
```

Enable SELinux permission for nginx proxying:

```bash
sudo setsebool -P httpd_can_network_connect on
```

Setup certbot:

<https://certbot.eff.org/instructions?ws=nginx&os=pip>

```bash
# I might have had to find an alternative way for augeas for the linux image I used
sudo dnf install python3 python-devel augeas gcc

sudo python3 -m venv /opt/certbot/
sudo /opt/certbot/bin/pip install --upgrade pip

sudo /opt/certbot/bin/pip install certbot certbot-nginx

sudo ln -s /opt/certbot/bin/certbot /usr/bin/certbot

sudo certbot certonly --nginx

echo "0 0,12 * * * root /opt/certbot/bin/python -c 'import random; import time; time.sleep(random.random() * 3600)' && sudo certbot renew -q" | sudo tee -a /etc/crontab > /dev/null
```

Reload nginx:

```bash
sudo nginx -t
sudo systemctl reload nginx
```
