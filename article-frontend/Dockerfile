FROM node:18-alpine as build

WORKDIR /app

ENV PATH /app/node_modules/.bin:$PATH

COPY . ./

RUN npm install

RUN npm run build

# production environment
FROM nginx:alpine

ARG ARTICLES_BACKEND_URL="http://localhost:1200/articles"

ENV ARTICLES_BACKEND_URL=${ARTICLES_BACKEND_URL}

COPY --from=build /app/build /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
