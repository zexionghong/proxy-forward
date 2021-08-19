rsync -avP -e "ssh" http-proxy-forward deploy@168.63.153.150:/home/deploy/vhost/http-proxy-forward/
rsync -avP -e "ssh" socks-proxy-forward deploy@168.63.153.150:/home/deploy/vhost/socks-proxy-forward/
rsync -avP -e "ssh" GeoLite2-City.mmdb deploy@168.63.153.150:/home/deploy/vhost/socks-proxy-forward/
