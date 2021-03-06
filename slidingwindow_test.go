package average

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	_, err := New(time.Second, time.Second)
	if err == nil || err.Error() != "window size has to be a multiplier of the granularity size" {
		t.Errorf("expected multiplier error, not %q", err)
	}
	_, err = New(time.Second, 2*time.Second)
	if err == nil || err.Error() != "window size has to be a multiplier of the granularity size" {
		t.Errorf("expected multiplier error, not %q", err)
	}
	_, err = New(3*time.Second, 2*time.Second)
	if err == nil || err.Error() != "window size has to be a multiplier of the granularity size" {
		t.Errorf("expected multiplier error, not %q", err)
	}

	_, err = New(0, time.Second)
	if err == nil || err.Error() != "window cannot be 0" {
		t.Errorf("expected window size cannot be 0 error, not %q", err)
	}

	_, err = New(time.Second, 0)
	if err == nil || err.Error() != "granularity cannot be 0" {
		t.Errorf("expected granularity cannot be 0 error, not %q", err)
	}
}

func TestAdd(t *testing.T) {
	sw := &SlidingWindow{
		window:      2 * time.Second,
		granularity: time.Second,
		samples:     []float64{1, 1},
		counts:      []int64{1, 2},
		pos:         1,
		size:        2,
	}

	sw.Add(1)
	if v := sw.samples[1]; v != 2 {
		t.Errorf("expected the 2nd sample to be 2, but got %f", v)
	}
}

func TestAverage(t *testing.T) {
	sw := &SlidingWindow{
		window:      10 * time.Second,
		granularity: time.Second,
		samples:     []float64{20, 4, 5, 0, 0, 0, 0, 0, 4, 10},
		counts:      []int64{10, 2, 5, 0, 0, 0, 0, 0, 4, 2},
		pos:         1,
		size:        10,
	}

	assert.Equal(t, 0.0, sw.Average(0))
	assert.Equal(t, 2.0, sw.Average(time.Second))
	assert.Equal(t, 2.0, sw.Average(2*time.Second))
	assert.Equal(t, 2.111111111111111, sw.Average(4*time.Second))
	assert.Equal(t, 1.8695652173913044, sw.Average(10*time.Second))
	assert.Equal(t, 1.8695652173913044, sw.Average(20*time.Second))
}

func TestReset(t *testing.T) {
	sw := MustNew(2*time.Second, time.Second)
	defer sw.Stop()

	sw.samples = []float64{1, 2}
	sw.pos = 1
	sw.size = 10

	sw.Reset()
	for _, v := range sw.samples {
		if v != 0 {
			t.Fatalf("expected the samples all to be 0, but at least one value was %f", v)
		}
	}
}

func TestResetFlow(t *testing.T) {
	sw := MustNew(time.Second, 10*time.Millisecond)
	defer sw.Stop()

	sw.Reset()
	time.Sleep(50 * time.Millisecond)
	sw.Reset()
	time.Sleep(50 * time.Millisecond)
	sw.Reset()
}

func TestTotal(t *testing.T) {
	sw := &SlidingWindow{
		window:      10 * time.Second,
		granularity: time.Second,
		samples:     []float64{1, 2, 5, 0, 0, 0, 0, 0, 4, 0},
		counts:      []int64{1, 2, 2, 0, 0, 0, 0, 0, 4, 0},
		pos:         1,
		size:        10,
	}

	if v, _ := sw.Total(0); v != 0 {
		t.Errorf("expected the total with a window of 0 seconds to be 0, not %f", v)
	}
	if v, _ := sw.Total(time.Second); v != 2 {
		t.Errorf("expected the total over the last second to be 2, not %f", v)
	}
	if v, _ := sw.Total(2 * time.Second); v != 3 {
		t.Errorf("expected the total over the last 2 seconds to be 3, not %f", v)
	}
	if v, _ := sw.Total(4 * time.Second); v != 7 {
		t.Errorf("expected the total over the last 4 seconds to be 7, not %f", v)
	}
	if v, _ := sw.Total(10 * time.Second); v != 12 {
		t.Errorf("expected the total over the last 10 seconds to be 12, not %f", v)
	}
	// This one should be equivalent to 10 seconds.
	if v, _ := sw.Total(20 * time.Second); v != 12 {
		t.Errorf("expected the total over the last 10 seconds to be 12, not %f", v)
	}
}

func TestTotalFromNew(t *testing.T) {
	sw := MustNew(10*time.Second, time.Second)

	sw.Add(10)
	sw.Add(20)
	total, samples := sw.Total(time.Second)

	assert.Equal(t, 30.0, total)
	assert.Equal(t, 1, samples)
}
