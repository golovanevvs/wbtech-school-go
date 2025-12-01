// Package or provides a utility for merging multiple done-channels into one.
package or

// Or merges multiple done-channels into a single resulting channel.
func Or(channels ...<-chan any) <-chan any {
	switch len(channels) {
	case 0:
		ch := make(chan any)
		close(ch)
		return ch
	case 1:
		return channels[0]
	}

	orDone := make(chan any)

	go func() {
		defer close(orDone)

		switch len(channels) {
		case 2:
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default:
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-channels[2]:
			case <-Or(append(channels[3:], orDone)...):
			}
		}
	}()

	return orDone
}
