# Build stage
FROM node:18-alpine as build

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm install --legacy-peer-deps

# Copy source code
COPY . .

# Build the application
RUN npm run build

# Production stage
FROM node:18-alpine

# timezone 패키지 설치
RUN apk --no-cache add tzdata

WORKDIR /app

# Copy built application from build stage
COPY --from=build /app/build /app/build

# Install a simple static file server
RUN npm install -g serve

# Expose port 3000
EXPOSE 3000

# Start the static file server
CMD ["serve", "-s", "build", "-l", "3000"]
