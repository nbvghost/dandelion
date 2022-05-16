docker pull ghcr.io/nbvghost/dandelion_gateway/server:latest
docker pull ghcr.io/nbvghost/oss/server:latest
docker pull ghcr.io/nbvghost/dandelion_site/server:latest
docker pull ghcr.io/nbvghost/dandelion_admin/server:latest
docker pull ghcr.io/nbvghost/sso/server:latest
docker-compose -f dandelion-docker-compose.yml down
docker-compose -f dandelion-docker-compose.yml up -d
