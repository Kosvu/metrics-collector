package retry

import "time"

func WithRetry(fn func() error, isRetriable func(error) bool) error {
	delays := []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}

	err := fn()

	if err == nil {
		return nil
	}

	for i := 0; i < 3; i++ {
		if !isRetriable(err) {
			return err
		}

		time.Sleep(delays[i])

		err = fn()

		if err == nil {
			return nil
		}
	}

	return err
}
