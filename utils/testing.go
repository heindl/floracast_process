package utils

import (
	"gopkg.in/mgo.v2"
	"github.com/saleswise/errors/errors"
	"github.com/mongodb/mongo-tools/mongoimport"
	"github.com/mongodb/mongo-tools/common/options"
	"github.com/facebookgo/mgotest"
	"github.com/mongodb/mongo-tools/common/db"
	"strconv"
)

func ImportFileToMongo(server *mgotest.Server, col *mgo.Collection, file string) error {
	tooloptions := options.New("test", "", options.EnabledOptions{false, false, false})
	tooloptions.Host = "localhost"
	tooloptions.Port = strconv.Itoa(server.Port)
	tooloptions.Timeout = 3
	tooloptions.Collection = col.Name
	tooloptions.DB = col.Database.Name
	sessionprovider, err := db.NewSessionProvider(*tooloptions)
	if err != nil {
		return errors.Wrap(err, "could not get session provider")
	}
	importer := mongoimport.MongoImport{
		ToolOptions: tooloptions,
		InputOptions: &mongoimport.InputOptions{
			File: file,
			HeaderLine: false,
			Type: mongoimport.JSON,
			JSONArray: true,
		},
		SessionProvider: sessionprovider,
		IngestOptions: &mongoimport.IngestOptions{},
	}
	if err := importer.ValidateSettings(nil); err != nil {
		return errors.Wrap(err, "could not validate settings")
	}
	if _, err := importer.ImportDocuments(); err != nil {
		return errors.Wrap(err, "could not import documents")
	}
	return nil
}
