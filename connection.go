package mongorm

import (
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

var connections = make(map[string]*mongo.Client)

func (m *MongORM[T]) initializeClient() {
	var client *mongo.Client
	if m.options != nil && m.options.MongoClient != nil {
		client = m.options.MongoClient
	} else {
		client = m.getClientFromSchema()
	}

	if client == nil {
		panic("MongoDB client connetion is not provided in options or schema")
	}

	if m.options != nil && m.options.DatabaseName != nil {
		m.info.db = client.Database(*m.options.DatabaseName)
	} else {
		m.setDatabaseFromSchema(client)
	}

	if m.info.db == nil {
		panic("MongoDB database is not provided in options or schema")
	}

	m.info.dbName = String(m.info.db.Name())

	if m.options != nil && m.options.CollectionName != nil {
		m.info.collection = m.info.db.Collection(*m.options.CollectionName)
	} else {
		m.setCollectionFromSchema()
	}

	if m.info.collection == nil {
		panic("MongoDB collection is not provided in options or schema")
	}

	m.setTimestampRequirementsFromSchema()
}

func getConnectionStringFromSchema[T any](schema *T) *string {
	ref := reflect.ValueOf(schema).Elem()
	t := ref.Type()

	for i := 0; i < ref.NumField(); i++ {
		fieldType := t.Field(i)

		// Skip exported
		if fieldType.PkgPath == "" {
			continue
		}

		if doesModelIncludeAnyModelFlags(fieldType.Tag, string(ModelTagConnectionString)) {
			tags := getModelTags(fieldType.Tag)
			if len(tags) <= 1 {
				panic(fmt.Sprintf("Field %s is missing the connection string tag value", fieldType.Name))
			}

			return String(tags[0])
		}
	}

	return nil
}

func (m *MongORM[T]) getClientFromSchema() *mongo.Client {
	connectionString := getConnectionStringFromSchema(m.schema)
	if connectionString == nil {
		panic("Connection string was not provided for database connection")
	}

	if connections[*connectionString] != nil {
		return connections[*connectionString]
	}

	client, err := NewClient(*connectionString)
	if err != nil {
		panic(fmt.Sprintf("Failed to create MongoDB client: %v", err))
	}

	connections[*connectionString] = client

	return client
}

func (m *MongORM[T]) setDatabaseFromSchema(client *mongo.Client) {
	ref := reflect.ValueOf(m.schema).Elem()
	t := ref.Type()

	for i := 0; i < ref.NumField(); i++ {
		fieldType := t.Field(i)

		// Skip exported
		if fieldType.PkgPath == "" {
			continue
		}

		if doesModelIncludeAnyModelFlags(fieldType.Tag, string(ModelTagDatabase)) {
			tags := getModelTags(fieldType.Tag)
			if len(tags) <= 1 {
				panic(fmt.Sprintf("Field %s is missing the database name tag value", fieldType.Name))
			}

			m.info.db = client.Database(tags[0])
			return
		}
	}
}

func (m *MongORM[T]) setCollectionFromSchema() {
	if m.info.db == nil {
		panic("Database is not configured")
	}

	ref := reflect.ValueOf(m.schema).Elem()
	t := ref.Type()

	for i := 0; i < ref.NumField(); i++ {
		fieldType := t.Field(i)

		// Skip exported
		if fieldType.PkgPath == "" {
			continue
		}

		if doesModelIncludeAnyModelFlags(fieldType.Tag, string(ModelTagCollection)) {
			tags := getModelTags(fieldType.Tag)
			if len(tags) <= 1 {
				panic(fmt.Sprintf("Field %s is missing the collection name tag value", fieldType.Name))
			}

			m.info.collection = m.info.db.Collection(tags[0])
			return
		}
	}
}
