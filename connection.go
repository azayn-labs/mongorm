package mongorm

import (
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// connections is a package-level variable that holds a map of MongoDB client connections.
// The keys of the map are connection strings, and the values are pointers to mongo.Client
// instances. This allows for reusing existing connections based on the connection string,
// improving performance and resource management by avoiding unnecessary creation of new
// clients for the same connection string. and resource management by avoiding unnecessary
// creation of new clients for the same connection string.
//
// > NOTE: This variable is internal only and should not be accessed
// directly from outside the package.
var connections = make(map[string]*mongo.Client)

func NewClient(connectionString string) (*mongo.Client, error) {
	client, err := mongo.Connect(
		options.Client().
			ApplyURI(connectionString).
			SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1)),
	)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (m *MongORM[T]) initializeClient() error {
	var client *mongo.Client
	var err error
	if m.options != nil && m.options.MongoClient != nil {
		client = m.options.MongoClient
	} else {
		client, err = m.getClientFromSchema()
		if err != nil {
			return err
		}
	}

	if client == nil {
		return fmt.Errorf("mongodb client connection is not provided in options or schema")
	}

	if m.options != nil && m.options.DatabaseName != nil {
		m.info.db = client.Database(*m.options.DatabaseName)
	} else {
		if err := m.setDatabaseFromSchema(client); err != nil {
			return err
		}
	}

	if m.info.db == nil {
		return fmt.Errorf("mongodb database is not provided in options or schema")
	}

	m.info.dbName = String(m.info.db.Name())

	if m.options != nil && m.options.CollectionName != nil {
		m.info.collection = m.info.db.Collection(*m.options.CollectionName)
	} else {
		if err := m.setCollectionFromSchema(); err != nil {
			return err
		}
	}

	if m.info.collection == nil {
		return fmt.Errorf("mongodb collection is not provided in options or schema")
	}

	if err := m.setTimestampRequirementsFromSchema(); err != nil {
		return err
	}

	return nil
}

func getConnectionStringFromSchema[T any](schema *T) (*string, error) {
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
				return nil, fmt.Errorf("field %s is missing the connection string tag value", fieldType.Name)
			}

			return String(tags[0]), nil
		}
	}

	return nil, nil
}

func (m *MongORM[T]) getClientFromSchema() (*mongo.Client, error) {
	connectionString, err := getConnectionStringFromSchema(m.schema)
	if err != nil {
		return nil, err
	}
	if connectionString == nil {
		return nil, fmt.Errorf("connection string was not provided for database connection")
	}

	if connections[*connectionString] != nil {
		return connections[*connectionString], nil
	}

	client, err := NewClient(*connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to create mongodb client: %w", err)
	}

	connections[*connectionString] = client

	return client, nil
}

func (m *MongORM[T]) setDatabaseFromSchema(client *mongo.Client) error {
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
				return fmt.Errorf("field %s is missing the database name tag value", fieldType.Name)
			}

			m.info.db = client.Database(tags[0])
			return nil
		}
	}

	return nil
}

func (m *MongORM[T]) setCollectionFromSchema() error {
	if m.info.db == nil {
		return fmt.Errorf("database is not configured")
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
				return fmt.Errorf("field %s is missing the collection name tag value", fieldType.Name)
			}

			m.info.collection = m.info.db.Collection(tags[0])
			return nil
		}
	}

	return nil
}
