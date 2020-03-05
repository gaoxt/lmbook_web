# ---- cache node_module ----
FROM node:14.4.0-alpine as base
COPY web-book/package.json .
RUN npm install -g cnpm --registry=https://registry.npm.taobao.org
RUN cnpm i
RUN cp -R node_modules prod_node_modules

# ---- build dist ----
FROM base AS dependencies
WORKDIR /app
COPY --from=base prod_node_modules node_modules
COPY web-book .
RUN npm run build

# ---- proxy dist----
FROM nginx:1.16.0-alpine 
WORKDIR /app
COPY --from=dependencies /app/dist .
COPY ./nginx.conf /etc/nginx/conf.d/default.conf
CMD ["nginx", "-g", "daemon off;"]