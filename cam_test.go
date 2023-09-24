package filtered-camera

import (
	"context"
	"image"
	"testing"

	"go.viam.com/rdk/resource"
	viz "go.viam.com/rdk/vision"
	"go.viam.com/rdk/vision/classification"
	"go.viam.com/rdk/vision/objectdetection"

	"go.viam.com/test"
)

type dummyVisionService struct {
	resource.AlwaysRebuild
	resource.TriviallyCloseable
	resource.Named
}

func (s *dummyVisionService) DetectionsFromCamera(ctx context.Context, cameraName string, extra map[string]interface{}) ([]objectdetection.Detection, error) {
	panic(1)
}

func (s *dummyVisionService) Detections(ctx context.Context, img image.Image, extra map[string]interface{}) ([]objectdetection.Detection, error) {

	if img == c {
		return []objectdetection.Detection{objectdetection.NewDetection(image.Rect(1, 1, 1, 1), .1, "b")}, nil
	}

	if img == b {
		return []objectdetection.Detection{objectdetection.NewDetection(image.Rect(1, 1, 1, 1), .9, "b")}, nil
	}

	return []objectdetection.Detection{}, nil
}

func (s *dummyVisionService) ClassificationsFromCamera(
	ctx context.Context,
	cameraName string,
	n int,
	extra map[string]interface{},
) (classification.Classifications, error) {
	panic(1)
}

func (s *dummyVisionService) Classifications(
	ctx context.Context,
	img image.Image,
	n int,
	extra map[string]interface{},
) (classification.Classifications, error) {

	if img == a {
		return classification.Classifications{classification.NewClassification(.9, "a")}, nil
	}

	if img == b {
		return classification.Classifications{classification.NewClassification(.1, "a")}, nil
	}

	return classification.Classifications{}, nil
}

func (s *dummyVisionService) GetObjectPointClouds(ctx context.Context, cameraName string, extra map[string]interface{}) ([]*viz.Object, error) {
	panic(1)
}

var (
	a = image.NewGray(image.Rect(1, 1, 1, 1))
	b = image.NewGray(image.Rect(2, 1, 1, 1))
	c = image.NewGray(image.Rect(3, 1, 1, 1))
	d = image.NewGray(image.Rect(4, 1, 1, 1))
)

func TestShouldSend(t *testing.T) {
	fc := &filteredCamera{
		conf: &Config{
			Classifications: map[string]float64{"a": .8},
			Objects:         map[string]float64{"b": .8},
		},
		vis: &dummyVisionService{},
	}

	res, err := fc.shouldSend(context.Background(), d)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, res, test.ShouldEqual, false)

	res, err = fc.shouldSend(context.Background(), c)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, res, test.ShouldEqual, false)

	res, err = fc.shouldSend(context.Background(), b)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, res, test.ShouldEqual, true)

	res, err = fc.shouldSend(context.Background(), a)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, res, test.ShouldEqual, true)
}
