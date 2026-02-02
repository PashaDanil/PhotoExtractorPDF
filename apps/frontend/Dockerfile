# syntax=docker/dockerfile:1

# Build stage
FROM node:25-alpine3.23 AS builder

# Set working directory
WORKDIR /build

# Copy package files first for better layer caching
COPY package.json package-lock.json* ./

# Install dependencies with clean install
RUN npm ci --only=production=false \
    && npm cache clean --force

# Copy source code
COPY . .

# Build the application
RUN npm run build

# Production stage with nginx
FROM nginx:1.28-alpine3.23

# Install runtime utilities
RUN apk add --no-cache curl

# Copy nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf

# Copy built application from builder
COPY --from=builder /build/dist /usr/share/nginx/html

# Set proper permissions
RUN chown -R nginx:nginx /usr/share/nginx/html && \
    chmod -R 755 /usr/share/nginx/html

# Expose port
EXPOSE 80

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost/health || exit 1

# Start nginx
CMD ["nginx", "-g", "daemon off;"]