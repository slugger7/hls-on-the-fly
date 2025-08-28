package m3u8

import (
	"fmt"
	"testing"

	"github.com/go-playground/assert"
)

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
	segments := generateSegmentsForManifest(5, nil, nameFunc)

	assert.IsEqual(segments, []Segment{})
}

func Test_generateSegmentForManifest_withNoFrames_shouldReturnEmptyArray(t *testing.T) {
	segments := generateSegmentsForManifest(5, []float64{}, nameFunc)

	assert.IsEqual(segments, []Segment{})
}

func Test_generateSegmentsForManifest_with_two_frames_one_more_than_hls_time_should_have_one_segment(t *testing.T) {
	hlsTime := 5
	duration := 5.1
	frames := []float64{0.0, duration}

	segments := generateSegmentsForManifest(hlsTime, frames, nameFunc)

	expected := []Segment{{Name: "0", Start: 0, Duration: duration}}

	assertEqualSegments(t, segments, expected)
}

func Test_generateSegmentsForManifest_with_two_frames_one_less_than_hls_time_should_have_one_segment(t *testing.T) {
	hlsTime := 5
	duration := 4.9
	frames := []float64{0.0, duration}

	segments := generateSegmentsForManifest(hlsTime, frames, nameFunc)

	expected := []Segment{{Name: "0", Start: 0, Duration: duration}}

	assertEqualSegments(t, segments, expected)
}

func Test_generateSegmentsForManifest_withThreeFrames_theLastWithinHlsTimeOfSecond_shouldReturnTwoSegments(t *testing.T) {
	hlsTime := 5
	frames := []float64{0.0, 5.1, 8.0}

	segments := generateSegmentsForManifest(hlsTime, frames, nameFunc)

	expected := []Segment{
		{Name: "0", Start: 0, Duration: 5.1},
		{Name: "1", Start: 5.1, Duration: frames[2] - frames[1]}, // if you do not use the actual values then it does not equal the same things
	}

	assertEqualSegments(t, segments, expected)
}

func Test_generateSegmentForManifest_withThreeFrames_theLastEqualingAnHLSSegment_shouldReturnTwoSegments(t *testing.T) {
	hlsTime := 5
	frames := []float64{0.0, 5.1, 10.0}

	segments := generateSegmentsForManifest(hlsTime, frames, nameFunc)

	expected := []Segment{
		{Name: "0", Start: 0, Duration: frames[1]},
		{Name: "1", Start: frames[1], Duration: frames[2] - frames[1]},
	}

	assertEqualSegments(t, segments, expected)
}
