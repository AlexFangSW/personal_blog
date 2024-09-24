FROM node:21-alpine

ARG BACKEND_BASE_URL

ENV BACKEND_BASE_URL=$BACKEND_BASE_URL
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
