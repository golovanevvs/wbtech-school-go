# build
FROM node:20-alpine AS builder
WORKDIR /app

RUN npm install -g pnpm

COPY package.json pnpm-lock.yaml ./

RUN pnpm install --prod --frozen-lockfile
RUN pnpm add typescript @types/node

COPY . .

ARG NEXT_PUBLIC_API_URL
ARG NEXT_PUBLIC_BASE_PATH
ENV NEXT_PUBLIC_API_URL=${NEXT_PUBLIC_API_URL}
ENV NEXT_PUBLIC_BASE_PATH=${NEXT_PUBLIC_BASE_PATH}

RUN pnpm run build

# prod
FROM node:20-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production

RUN npm install -g pnpm

COPY --from=builder /app/package.json ./
COPY --from=builder /app/pnpm-lock.yaml ./
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
COPY --from=builder /app/tsconfig.json ./
COPY --from=builder /app/next.config.ts ./

RUN pnpm install --prod --frozen-lockfile

# start
CMD ["pnpm", "start"]
