FROM node:16-alpine

# testting
# Install nano and curl
RUN apk --no-cache add nano curl


WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

RUN npm run build

EXPOSE 4000

CMD ["node", "dist/server.js"]