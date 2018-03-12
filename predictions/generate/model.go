package generate

import (
	"bitbucket.org/heindl/process/nameusage/nameusage"
	"bitbucket.org/heindl/process/store"
	"context"
	"fmt"
	"github.com/dropbox/godropbox/errors"
	tg "github.com/galeone/tfgo"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

type Modeller interface {
	FetchModel(ctx context.Context, usageID nameusage.ID) (*tg.Model, error)
	Close() error
}

func NewModeller(floraStore store.FloraStore) (Modeller, error) {
	return &modeller{
		floraStore:       floraStore,
		modelDirectories: []string{},
	}, nil
}

type modeller struct {
	floraStore       store.FloraStore
	modelDirectories []string
}

func (Ω *modeller) Close() error {
	// Cleanup old files
	for _, dirName := range Ω.modelDirectories {
		if !strings.HasPrefix(dirName, "/tmp") {
			return errors.Newf("Expected model directory [%s] to be in the /tmp folder ", dirName)
		}
		if err := os.RemoveAll(dirName); err != nil {
			return errors.Wrap(err, "Could not remove model directory")
		}
	}
	return nil
}

func (Ω *modeller) FetchModel(ctx context.Context, usageID nameusage.ID) (*tg.Model, error) {

	latestDate, err := Ω.fetchLatestModelDate(ctx, usageID)
	if err != nil {
		return nil, err
	}

	cloudPath := fmt.Sprintf("models/%s/%s", usageID.String(), latestDate)
	tempModelPath, err := ioutil.TempDir("/tmp", usageID.String()+"_")
	if err != nil {
		return nil, errors.Wrap(err, "Could not generate temporary directory for tensorflow model")
	}
	fmt.Println("MODEL_PATH", tempModelPath)
	Ω.modelDirectories = append(Ω.modelDirectories, tempModelPath)

	if err := Ω.floraStore.SyncGCSPathWithLocal(ctx, cloudPath, tempModelPath); err != nil {
		return nil, err
	}

	return tg.LoadModel(tempModelPath, []string{"serve"}, nil), nil
}

func (Ω *modeller) fetchLatestModelDate(ctx context.Context, usageID nameusage.ID) (string, error) {
	gcsPath := fmt.Sprintf("models/%s/", usageID.String())
	pbObjects, err := Ω.floraStore.CloudStorageObjects(ctx, gcsPath, ".pb")
	if err != nil {
		return "", err
	}
	dates := []string{}
	for _, pbObject := range pbObjects {
		attrs, err := pbObject.Attrs(ctx)
		if err != nil {
			return "", errors.Wrap(err, "Could not get object attributes")
		}
		d := strings.Split(strings.TrimPrefix(attrs.Name, gcsPath), "/")[0]
		dates = append(dates, d)
	}
	sort.Strings(dates)
	return dates[len(dates)-1], nil
}
