FROM node:20.11.1-alpine3.18

WORKDIR /sunflower

COPY . .
RUN npm ci --legacy-peer-deps
RUN npm run build
RUN rm -rf ./src
RUN rm -rf ./node_modules
RUN npm install -g serve
EXPOSE 3000
