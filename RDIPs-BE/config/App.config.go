package config

import (
	"context"
	"database/sql"
	"os"
	"slices"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	LogConstant "RDIPs-BE/constant/LogConst"
	AMQP_handler "RDIPs-BE/handler/AMQP"
	Keycloak_handler "RDIPs-BE/handler/Keycloak"

	"RDIPs-BE/model"
	"RDIPs-BE/utils"
)

type objDB interface {
	TableName() string
}

func MongoConfig() (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*60*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().
		ApplyURI("mongodb://"+os.Getenv("MONGO_HOST")+":27017").
		SetAuth(options.Credential{
			AuthMechanism: "SCRAM-SHA-256",
			AuthSource:    "admin",
			Username:      os.Getenv("MONGO_INITDB_ROOT_USERNAME"),
			Password:      os.Getenv("MONGO_INITDB_ROOT_PASSWORD"),
		}))
	if err != nil {
		return nil, err
	}

	mongoDB := client.Database(os.Getenv("MONGO_INITDB_DATABASE"))

	// Create all collection if not exist
	isNameOnly := true
	collections := []map[string]string{
		{
			"name":      "performance",
			"metaField": "document_name",
		},
	}
	currentCollection, err := mongoDB.ListCollectionNames(ctx, nil, &options.ListCollectionsOptions{
		NameOnly: &isNameOnly,
	})

	if err != nil && err != mongo.ErrNilDocument {
		return nil, err
	}

	for _, collection := range collections {
		if !slices.Contains(currentCollection, collection["name"]) {
			metaField := collection["metaField"]
			err = mongoDB.CreateCollection(ctx, collection["name"], &options.CreateCollectionOptions{
				TimeSeriesOptions: &options.TimeSeriesOptions{
					TimeField: "timestamp",
					MetaField: &metaField,
				},
			})
			if err != nil {
				return nil, err
			}
		}
	}

	return mongoDB, err
}

func DbConfig() (*gorm.DB, error) {
	sql := &sql.DB{}
	sql.SetMaxIdleConns(10)
	sql.SetConnMaxLifetime(1 * time.Minute)
	sql.SetConnMaxIdleTime(time.Second * 30)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: "host=" + os.Getenv("POSTGRES_HOST") +
			" user=" + os.Getenv("POSTGRES_USER") +
			" password=" + os.Getenv("POSTGRES_PASSWORD") +
			" dbname=" + os.Getenv("POSTGRES_DB") +
			" port=" + os.Getenv("POSTGRES_PORT") +
			" sslmode=disable TimeZone=Asia/Ho_Chi_Minh",
	}), &gorm.Config{
		ConnPool: sql,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: os.Getenv("POSTGRES_SCHEMA") + ".", // schema name
		}})
	if err != nil {
		return nil, err
	}

	err = db.Exec("CREATE SCHEMA IF NOT EXISTS " + os.Getenv("POSTGRES_SCHEMA")).Error

	if err == nil {
		models := []objDB{model.SysDevices{}, model.SysHistory{}, model.SysPerformance{}}
		relModels := []objDB{model.SysDeviceRel{}}
		for _, m := range append(models, relModels...) {
			if !db.Migrator().HasTable(m) {
				utils.Log(LogConstant.Info, "Create table "+m.TableName())
				err := db.Migrator().CreateTable(m)
				if err != nil {
					utils.Log(LogConstant.Warning, err)
				}
			}
		}
	}
	return db, err
}

func RabbitMqConfig() error {
	return AMQP_handler.InitializeAMQP()
}

func KeycloakConfig() error {
	return Keycloak_handler.InitKeycloakClient("")
}
