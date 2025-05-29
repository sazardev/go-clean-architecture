# mongodb/ - Implementación MongoDB

Implementaciones de repositorios usando MongoDB para casos de uso de datos no relacionales.

## Responsabilidades

- Implementar repositorios para datos no estructurados
- Manejar documentos flexibles
- Optimizar agregaciones complejas
- Gestionar colecciones y índices
- Implementar búsquedas de texto completo

## Casos de Uso

MongoDB es ideal para:
- **Logs y auditoría**: Documentos de eventos
- **Documentos flexibles**: Datos que cambian estructura
- **Agregaciones complejas**: Reportes y análisis
- **Datos jerárquicos**: Estructuras anidadas
- **Prototipos rápidos**: Desarrollo ágil

## Archivos Futuros

- **`audit_repository.go`** - Logs de auditoría
- **`document_repository.go`** - Documentos flexibles
- **`analytics_repository.go`** - Datos analíticos
- **`models.go`** - Modelos de MongoDB
- **`indexes.go`** - Definición de índices

## Características MongoDB

### Ventajas
- Esquema flexible
- Escalabilidad horizontal
- Agregaciones potentes
- Búsqueda de texto nativo
- Replicación automática

### Funcionalidades
- **Collections**: Equivalente a tablas
- **Documents**: Documentos BSON flexibles
- **Aggregation Pipeline**: Para análisis complejos
- **GridFS**: Para archivos grandes
- **Change Streams**: Para eventos en tiempo real

## Implementación

```go
type MongoAuditRepository struct {
    collection *mongo.Collection
}

func NewMongoAuditRepository(db *mongo.Database) *MongoAuditRepository {
    return &MongoAuditRepository{
        collection: db.Collection("audit_logs"),
    }
}

type AuditLog struct {
    ID        primitive.ObjectID `bson:"_id,omitempty"`
    UserID    string            `bson:"user_id"`
    Action    string            `bson:"action"`
    Resource  string            `bson:"resource"`
    Details   map[string]interface{} `bson:"details"`
    Timestamp time.Time         `bson:"timestamp"`
}

func (r *MongoAuditRepository) CreateLog(ctx context.Context, log *AuditLog) error {
    log.ID = primitive.NewObjectID()
    log.Timestamp = time.Now()
    
    _, err := r.collection.InsertOne(ctx, log)
    return err
}

func (r *MongoAuditRepository) FindLogsByUser(ctx context.Context, userID string, limit int) ([]*AuditLog, error) {
    filter := bson.M{"user_id": userID}
    opts := options.Find().SetLimit(int64(limit)).SetSort(bson.D{{"timestamp", -1}})
    
    cursor, err := r.collection.Find(ctx, filter, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)
    
    var logs []*AuditLog
    if err := cursor.All(ctx, &logs); err != nil {
        return nil, err
    }
    
    return logs, nil
}
```

## Agregaciones

```go
func (r *MongoAuditRepository) GetUserActivityStats(ctx context.Context, userID string) (*ActivityStats, error) {
    pipeline := mongo.Pipeline{
        bson.D{{"$match", bson.D{{"user_id", userID}}}},
        bson.D{{"$group", bson.D{
            {"_id", "$action"},
            {"count", bson.D{{"$sum", 1}}},
        }}},
        bson.D{{"$sort", bson.D{{"count", -1}}}},
    }
    
    cursor, err := r.collection.Aggregate(ctx, pipeline)
    // ... procesar resultados
}
```

## Configuración

Variables de entorno:
- `MONGODB_URI` - Connection string completo
- `MONGODB_DATABASE` - Nombre de la base de datos
- `MONGODB_MAX_POOL_SIZE` - Tamaño del pool
- `MONGODB_TIMEOUT` - Timeout de conexión
