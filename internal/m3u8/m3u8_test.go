package m3u8

import (
	"fmt"
	"testing"

	"github.com/go-playground/assert"
)

const HLS_TIME = 5

func assertEqualSegments(t *testing.T, actual, expected []Segment) {
	assert.Equal(t, len(actual), len(expected))

	for i, segment := range actual {
		assert.Equal(t, segment.Name, expected[i].Name)
		assert.Equal(t, segment.Start, expected[i].Start)
		assert.Equal(t, segment.Duration, expected[i].Duration)
	}
}

func nameFunc(i int) string {
	return fmt.Sprintf("%d", i)
}

func Test_generateSegmentForManifest_withNilFrames_shouldReturnEmptyArray(t *testing.T) {
	segments := generateSegmentsForManifest(5, 0.0, nil, nameFunc)

	assert.IsEqual(segments, []Segment{})
}

func Test_generateSegmentForManifest_withNoFrames_shouldReturnEmptyArray(t *testing.T) {
	segments := generateSegmentsForManifest(5, 0.0, []float64{}, nameFunc)

	assert.IsEqual(segments, []Segment{})
}

func Test_generateSegmentForManifest_withOneFrameLessThanHLSTime_shouldReturnOneSegment(t *testing.T) {
	frames := []float64{0.0, 2.2}

	segments := generateSegmentsForManifest(HLS_TIME, 2.2, frames, nameFunc)

	expected := []Segment{{Name: nameFunc(0), Start: frames[0], Duration: frames[1]}}

	assertEqualSegments(t, segments, expected)
}

func Test_generateSegmentForManifest_withTwoFramesLessThanHLSTime_shouldReturnOneSegment(t *testing.T) {
	frames := []float64{0.0, 2.2, 4.4}

	segments := generateSegmentsForManifest(HLS_TIME, 4.4, frames, nameFunc)

	expected := []Segment{{Name: nameFunc(0), Start: frames[0], Duration: frames[2]}}

	assertEqualSegments(t, segments, expected)
}

func Test_generateSegmentForManifest_withOneFrameLessThanHLSTimeAndOneFrameAboveeHLSTime_shouldReturnTwoSegments(t *testing.T) {
	frames := []float64{0.0, 2.2, 5.1}

	segments := generateSegmentsForManifest(HLS_TIME, 5.1, frames, nameFunc)

	fmt.Println(segments)

	expected := []Segment{
		{Name: nameFunc(0), Start: frames[0], Duration: frames[1]},
		{Name: nameFunc(1), Start: frames[1], Duration: frames[2] - frames[1]},
	}

	assertEqualSegments(t, segments, expected)
}

func Test_generatSegmentForManifest_withFirstFrameUnderHLSTimeAndSecondFrameEqualToHLSTime_shouldReturnTwoSegmentsWhereSecondSegmentEqualsHLSTime(t *testing.T) {
	frames := []float64{0.0, 2.5, 7.5}

	segments := generateSegmentsForManifest(HLS_TIME, 7.5, frames, nameFunc)

	expected := []Segment{
		{Name: nameFunc(0), Start: frames[0], Duration: frames[1]},
		{Name: nameFunc(1), Start: frames[1], Duration: frames[2] - frames[1]},
	}

	assertEqualSegments(t, segments, expected)
}

func Test_generateSegmentForManifest_withLastFrameBetweenSegments_shouldReturnTwoSegments(t *testing.T) {
	frames := []float64{0.0, 5.0, 7.5}

	segments := generateSegmentsForManifest(HLS_TIME, 7.5, frames, nameFunc)

	expeted := []Segment{
		{Name: nameFunc(0), Start: frames[0], Duration: frames[1]},
		{Name: nameFunc(1), Start: frames[1], Duration: frames[2] - frames[1]},
	}

	assertEqualSegments(t, segments, expeted)
}
