#!/bin/bash


GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Starting e-commerce services...${NC}"

if ! command -v docker &> /dev/null; then
    echo -e "${YELLOW}Docker is not installed. Please install Docker first.${NC}"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo -e "${YELLOW}Docker Compose is not installed. Please install Docker Compose first.${NC}"
    exit 1
fi

echo -e "${GREEN}Building and starting services...${NC}"
docker-compose up -d

echo -e "${GREEN}Service status:${NC}"
docker-compose ps

echo -e "${GREEN}Services started! The API Gateway is available at: http://localhost:8080${NC}"
echo -e "${YELLOW}Note: You need to create the databases manually. See README.md for instructions.${NC}"
echo -e "${GREEN}You can check logs with: docker-compose logs -f${NC}" 