# ---- Base ----
# Use a specific Node.js version for reproducibility
FROM node:20-alpine AS base
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
# Enable pnpm, which will use the version specified in package.json
RUN corepack enable
WORKDIR /app

# ---- Dependencies ----
# This layer installs ALL dependencies (including dev) needed for the build.
# It is cached as long as package files don't change.
FROM base AS deps
COPY package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

# ---- Builder ----
# This stage builds the SvelteKit application.
FROM base AS builder
ARG GO_SERVER_ADDRESS
ARG WEBSOCKET_ADDRESS
ENV GO_SERVER_ADDRESS=${GO_SERVER_ADDRESS}
ENV WEBSOCKET_ADDRESS=${WEBSOCKET_ADDRESS}

# Copy source code and the installed dependencies
COPY . .
COPY --from=deps /app/node_modules ./node_modules
RUN pnpm run build

# ---- Production ----
# This is the final, small production image.
FROM base AS production
ENV NODE_ENV=production

# Copy package files to install only production dependencies
COPY package.json pnpm-lock.yaml ./
RUN pnpm install --prod --frozen-lockfile

# Copy the built application artifacts from the builder stage
COPY --from=builder /app/build ./build

EXPOSE 3000
CMD [ "node", "build" ]
