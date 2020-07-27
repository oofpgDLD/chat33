docker run \
-p 8090:8090 \
--name chat33-server \
--network chatnet \
--network-alias server \
-d chat33:1.0.0