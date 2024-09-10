package uticker

import (
	"testing"
	"time"
)

func TestNewUTicker(t *testing.T) {
	ticker := NewUTicker()
	defer ticker.Stop()

	if ticker.duration != time.Second {
		t.Errorf("Expected default duration to be 1 second, got %v", ticker.duration)
	}

	if ticker.immediateStart {
		t.Error("Expected immediateStart to be false by default")
	}

	if ticker.C == nil {
		t.Error("Expected channel C to be initialized")
	}
}

func TestWithImmediateStart(t *testing.T) {
	ticker := NewUTicker(WithImmediateStart())
	defer ticker.Stop()

	if !ticker.immediateStart {
		t.Error("Expected immediateStart to be true")
	}

	select {
	case <-ticker.C:
		// Immediate tick received
	case <-time.After(100 * time.Millisecond):
		t.Error("Expected immediate tick, but didn't receive one")
	}
}

func TestWithDuration(t *testing.T) {
	duration := 500 * time.Millisecond
	ticker := NewUTicker(WithDuration(duration))
	defer ticker.Stop()

	if ticker.duration != duration {
		t.Errorf("Expected duration to be %v, got %v", duration, ticker.duration)
	}

	start := time.Now()
	<-ticker.C
	elapsed := time.Since(start)

	if elapsed < duration {
		t.Errorf("Tick occurred too soon. Expected at least %v, got %v", duration, elapsed)
	}
}

func TestWithExponentialBackoff(t *testing.T) {
	ticker := NewUTicker(
		WithDuration(100*time.Millisecond),
		WithExponentialBackoff(2),
	)
	defer ticker.Stop()

	expectedDurations := []time.Duration{
		100 * time.Millisecond,
		200 * time.Millisecond,
		400 * time.Millisecond,
	}

	for i, expected := range expectedDurations {
		start := time.Now()
		<-ticker.C
		elapsed := time.Since(start)

		if elapsed < expected-50*time.Millisecond || elapsed > expected+50*time.Millisecond {
			t.Errorf("Tick %d: Expected duration close to %v, got %v", i, expected, elapsed)
		}
	}
}

func TestWithExponentialBackoffCapped(t *testing.T) {
	ticker := NewUTicker(
		WithDuration(100*time.Millisecond),
		WithExponentialBackoffCapped(2, 2),
	)
	defer ticker.Stop()

	expectedDurations := []time.Duration{
		100 * time.Millisecond,
		200 * time.Millisecond,
		400 * time.Millisecond,
		400 * time.Millisecond,
	}

	for i, expected := range expectedDurations {
		start := time.Now()
		<-ticker.C
		elapsed := time.Since(start)

		if elapsed < expected-50*time.Millisecond || elapsed > expected+50*time.Millisecond {
			t.Errorf("Tick %d: Expected duration close to %v, got %v", i, expected, elapsed)
		}
	}
}

func TestWithRampCapped(t *testing.T) {
	ticker := NewUTicker(
		WithDuration(400*time.Millisecond),
		WithRampCapped(2, 2),
	)
	defer ticker.Stop()

	expectedDurations := []time.Duration{
		400 * time.Millisecond,
		200 * time.Millisecond,
		100 * time.Millisecond,
		100 * time.Millisecond,
	}

	for i, expected := range expectedDurations {
		start := time.Now()
		<-ticker.C
		elapsed := time.Since(start)

		if elapsed < expected-50*time.Millisecond || elapsed > expected+50*time.Millisecond {
			t.Errorf("Tick %d: Expected duration close to %v, got %v", i, expected, elapsed)
		}
	}
}

func TestWithDeviation(t *testing.T) {
	baseDuration := 100 * time.Millisecond
	deviation := 0.5 // 50% deviation
	ticker := NewUTicker(
		WithDuration(baseDuration),
		WithDeviation(deviation),
	)
	defer ticker.Stop()

	for i := 0; i < 10; i++ {
		start := time.Now()
		<-ticker.C
		elapsed := time.Since(start)

		if elapsed < baseDuration || elapsed > baseDuration+baseDuration*time.Duration(deviation) {
			t.Errorf("Tick %d: duration %v outside expected range [%v, %v]", i, elapsed, baseDuration, baseDuration+baseDuration*time.Duration(deviation))
		}
	}
}

func TestWithAnotherDurationWithGivenProbability(t *testing.T) {
	baseDuration := 100 * time.Millisecond
	anotherDuration := 300 * time.Millisecond
	probability := 0.5
	ticker := NewUTicker(
		WithDuration(baseDuration),
		WithAnotherDurationWithGivenProbability(anotherDuration, probability),
	)
	defer ticker.Stop()

	baseDurationCount := 0
	anotherDurationCount := 0

	for i := 0; i < 100; i++ {
		start := time.Now()
		<-ticker.C
		elapsed := time.Since(start)

		if elapsed >= baseDuration-10*time.Millisecond && elapsed <= baseDuration+10*time.Millisecond {
			baseDurationCount++
		} else if elapsed >= anotherDuration-10*time.Millisecond && elapsed <= anotherDuration+10*time.Millisecond {
			anotherDurationCount++
		} else {
			t.Errorf("Tick %d: Unexpected duration %v", i, elapsed)
		}
	}

	if baseDurationCount < 30 || baseDurationCount > 70 || anotherDurationCount < 30 || anotherDurationCount > 70 {
		t.Errorf("Expected roughly equal distribution. Got base: %d, another: %d", baseDurationCount, anotherDurationCount)
	}
}

func TestWithRandomTickIn(t *testing.T) {
	maxDuration := 500 * time.Millisecond
	ticker := NewUTicker(
		WithRandomTickIn(maxDuration),
	)
	defer ticker.Stop()

	for i := 0; i < 10; i++ {
		start := time.Now()
		<-ticker.C
		elapsed := time.Since(start)

		if elapsed > maxDuration {
			t.Errorf("Tick %d: duration %v exceeded maximum %v", i, elapsed, maxDuration)
		}
	}
}

func TestStop(t *testing.T) {
	ticker := NewUTicker()
	ticker.Stop()

	select {
	case _, ok := <-ticker.C:
		if ok {
			t.Error("Channel should be closed after Stop()")
		}
	default:
		t.Error("Channel should be closed and not blocking")
	}
}

func TestReset(t *testing.T) {
	ticker := NewUTicker(WithDuration(100 * time.Millisecond))
	defer ticker.Stop()

	<-ticker.C // Wait for first tick

	newDuration := 200 * time.Millisecond
	ticker.Reset(newDuration)

	start := time.Now()
	<-ticker.C
	elapsed := time.Since(start)

	if elapsed < newDuration-50*time.Millisecond || elapsed > newDuration+50*time.Millisecond {
		t.Errorf("After Reset: Expected duration close to %v, got %v", newDuration, elapsed)
	}
}
