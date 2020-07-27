docker run \
-p 3308:3306 \
--name chat33-mysql \
--network chatnet \
--network-alias mysql \
-e MYSQL_ROOT_PASSWORD=123456 \
-d mysql:chat \
--bind-address=0.0.0.0