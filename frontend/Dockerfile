# ---- Base ----
# Use a specific Node.js version for reproducibility
FROM node:20-alpine AS base
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
# Enable pnpm
RUN corepack enable
WORKDIR /app

# ---- Dependencies ----
# This layer is cached as long as package files don't change
FROM base AS deps
COPY package.json pnpm-lock.yaml ./
RUN pnpm install --prod --frozen-lockfile

# ---- Builder ----
# This stage builds the application
FROM base AS builder
# Receive build-time environment variables from Docker Compose
ARG GO_SERVER_ADDRESS
ARG WEBSOCKET_ADDRESS

# Set them as environment variables for the SvelteKit build process
ENV GO_SERVER_ADDRESS=${GO_SERVER_ADDRESS}
ENV WEBSOCKET_ADDRESS=${WEBSOCKET_ADDRESS}

COPY . .
COPY --from=deps /app/node_modules ./node_modules

# Build the SvelteKit application
RUN pnpm run build

# ---- Production ----
# This is the final, small production image
FROM base AS production
ENV NODE_ENV=production

# Copy only the necessary build artifacts and production dependencies
COPY --from=builder /app/build ./build
COPY --from=builder /app/node_modules ./node_modules
COPY package.json .

# The @sveltejs/adapter-node produces a standalone Node.js server
EXPOSE 3000

# The "start" script in package.json is "node build"
CMD [ "node", "build" ]
