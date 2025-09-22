# Redis Implementation Documentation

## Overview
Redis caching telah diimplementasikan untuk brand operations untuk meningkatkan performance aplikasi dengan mengurangi query database yang berulang.

## Configuration (.env)

```properties
# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
REDIS_ENABLED=true
REDIS_TTL=3600
```

### Configuration Options:
- `REDIS_HOST`: Redis server host (default: localhost)
- `REDIS_PORT`: Redis server port (default: 6379)
- `REDIS_PASSWORD`: Redis password (optional)
- `REDIS_DB`: Redis database number (default: 0)
- `REDIS_ENABLED`: Enable/disable Redis caching (true/false)
- `REDIS_TTL`: Cache time-to-live in seconds (default: 3600 = 1 hour)

## Cache Keys Strategy

### Brand Operations:
- `brands:all` - Cache untuk GetAllBrands()
- `brand:id:{id}` - Cache untuk GetBrandByID(id)

## Implementation Details

### 1. Redis Module (`pkg/redis/redis.go`)
- Connection management
- Basic operations (Set, Get, Delete, DeletePattern)
- Environment-based configuration
- Graceful fallback when Redis unavailable

### 2. Brand Repository (`internal/repository/brand_repository.go`)
- Cache-first read operations
- Cache invalidation on write operations
- JSON serialization for complex data structures

### 3. Cache Strategies

#### Read Operations (Cache-First):
1. Check Redis cache first
2. If cache miss, query database
3. Store result in cache with TTL
4. Return data

#### Write Operations (Cache Invalidation):
1. Perform database operation
2. Invalidate related cache entries
3. Cache will be repopulated on next read

## Performance Benefits

### Before Redis:
- Every request hits database
- Repeated queries for same data
- Higher database load

### After Redis:
- Cache hit = instant response
- Reduced database queries
- Better scalability
- Improved response times

## Usage Examples

### Installing Redis (Local Development):

#### Windows (using Chocolatey):
```bash
choco install redis-64
redis-server
```

#### Docker:
```bash
docker run -d --name redis -p 6379:6379 redis:latest
```

#### Ubuntu/Linux:
```bash
sudo apt update
sudo apt install redis-server
sudo systemctl start redis-server
```

### Testing Redis Connection:
```bash
redis-cli ping
# Should return: PONG
```

## Cache Management

### Manual Cache Clearing:
```bash
# Clear specific brand cache
redis-cli DEL "brand:id:1"

# Clear all brands cache
redis-cli DEL "brands:all"

# Clear all cache
redis-cli FLUSHALL
```

### Monitoring Cache:
```bash
# View all keys
redis-cli KEYS "*"

# Monitor Redis operations
redis-cli MONITOR

# Check memory usage
redis-cli INFO memory
```

## Error Handling

Redis implementation includes graceful fallback:
- If Redis is disabled (`REDIS_ENABLED=false`), operations continue without caching
- If Redis connection fails, application logs warning but continues normally
- Cache errors are logged but don't break application flow

## Best Practices

1. **TTL Management**: Set appropriate TTL based on data change frequency
2. **Cache Invalidation**: Always invalidate cache on data modifications
3. **Memory Management**: Monitor Redis memory usage in production
4. **Key Naming**: Use consistent, descriptive key patterns
5. **Error Handling**: Implement graceful fallback for cache failures

## Production Considerations

1. **Redis Server Setup**: Use dedicated Redis server in production
2. **Security**: Configure Redis authentication and network security
3. **Monitoring**: Set up Redis monitoring and alerting
4. **Backup**: Configure Redis persistence (RDB/AOF)
5. **Scaling**: Consider Redis Cluster for high availability

## Future Enhancements

1. **Category Caching**: Implement similar caching for categories
2. **Product Caching**: Add caching for product operations
3. **Cache Warming**: Pre-populate cache for frequently accessed data
4. **Cache Analytics**: Track cache hit/miss rates
5. **Advanced Patterns**: Implement cache-aside, write-through patterns

## Troubleshooting

### Common Issues:

1. **Connection Refused**: Check if Redis server is running
2. **Authentication Failed**: Verify REDIS_PASSWORD in .env
3. **Cache Not Working**: Check REDIS_ENABLED=true in .env
4. **Memory Issues**: Monitor Redis memory usage and set maxmemory

### Debug Commands:
```bash
# Check Redis status
redis-cli info

# Test connection
redis-cli ping

# View current keys
redis-cli keys "*"

# Check specific key
redis-cli get "brands:all"
```