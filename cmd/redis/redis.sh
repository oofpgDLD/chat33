docker run \
-p 6380:6379 \
--name chat33-redis \
--network chatnet \
--network-alias redis \
-d redis:chat