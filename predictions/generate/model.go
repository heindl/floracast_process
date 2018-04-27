package generate

import (
	"bitbucket.org/heindl/process/nameusage/nameusage"
	"bitbucket.org/heindl/process/store"
	"context"
	"fmt"
	"github.com/dropbox/godropbox/errors"
	tg "github.com/galeone/tfgo"
	"github.com/golang/glog"
	"io/ioutil"
	"sort"
	"strings"
)

type FloraClassifier interface {
	NewClassifierInstance(ctx context.Context) (*tg.Model, error)
	Close() error
}

func NewFloraClassifier(cxt context.Context, floraStore store.FloraStore, nameUsageID nameusage.ID, modelPath string) (FloraClassifier, error) {
	c := &classifier{
		floraStore:  floraStore,
		nameUsageID: nameUsageID,
	}
	if modelPath != "" {
		c.modelPath = modelPath
	} else {
		if err := c.download(cxt); err != nil {
			return nil, err
		}
	}

	return c, nil

}

type classifier struct {
	floraStore       store.FloraStore
	modelDirectories []string
	nameUsageID      nameusage.ID
	modelPath        string
}

func (Ω *classifier) Close() error {
	// Cleanup old files
	//for _, dirName := range Ω.modelDirectories {
	//	if !strings.HasPrefix(dirName, "/tmp") {
	//		return errors.Newf("Expected model directory [%s] to be in the /tmp folder ", dirName)
	//	}
	//	if err := os.RemoveAll(dirName); err != nil {
	//		return errors.Wrap(err, "Could not remove model directory")
	//	}
	//}
	return nil
}

func (Ω *classifier) download(ctx context.Context) error {
	latestDate, err := Ω.fetchLatestModelDate(ctx, Ω.nameUsageID)
	if err != nil {
		return err
	}

	cloudPath := fmt.Sprintf("models/%s/%s", Ω.nameUsageID.String(), latestDate)
	tempModelPath, err := ioutil.TempDir("/tmp", Ω.nameUsageID.String()+"_")
	if err != nil {
		return errors.Wrap(err, "Could not generate temporary directory for tensorflow model")
	}

	glog.Infof("Downloading Floracast Classifier [%s] from GCS [%s] to Local Path [%s]", Ω.nameUsageID, cloudPath, tempModelPath)

	Ω.modelPath = tempModelPath

	return Ω.floraStore.SyncGCSPathWithLocal(ctx, cloudPath, tempModelPath)
}

func (Ω *classifier) NewClassifierInstance(ctx context.Context) (*tg.Model, error) {

	return tg.LoadModel(Ω.modelPath, []string{"serve"}, nil), nil
}

func (Ω *classifier) fetchLatestModelDate(ctx context.Context, usageID nameusage.ID) (string, error) {
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
