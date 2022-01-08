# Test create typesense
- `$ mkdir /tmp/typesense-data
docker run --name typesensedock -p 8108:8108 -v/tmp/data:/data typesense/typesense:0.21.0 --enable-cors --data-dir /data --api-key=Hu52dwsas2AdxdE`