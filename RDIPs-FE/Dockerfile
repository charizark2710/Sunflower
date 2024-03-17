FROM node:latest

WORKDIR /sunflower

COPY . .
ENV NODE_ENV production
RUN npm install -g serve
RUN npm ci
RUN npm run build
RUN rm -rf ./src
RUN rm -rf ./node_modules
EXPOSE 3000

CMD [ "serve", "-s", "build" ]


