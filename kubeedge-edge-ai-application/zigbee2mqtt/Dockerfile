FROM node:12 as base

WORKDIR /app


COPY package.json ./

# Dependencies
FROM base as dependencies

COPY npm-shrinkwrap.json ./

Run npm ci --production

# Release
FROM base as release

COPY --from=dependencies /app/node_modules ./node_modules
COPY lib ./lib
COPY index.js ./
COPY data/ ./data

COPY docker/docker-entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/docker-entrypoint.sh


ARG COMMIT
RUN echo "{\"hash\": \"$COMMIT\"}" > .hash.json

ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["npm", "start"]
