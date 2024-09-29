FROM node:21-alpine

ENV NODE_ENV=production

WORKDIR /app

COPY package* .env /app/

RUN npm install --include dev

COPY . .

RUN npm run build

EXPOSE 3000

RUN chown -R node:node /app/ 

USER node

CMD npm run start
