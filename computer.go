package main

// Computer is an interface for quantile computation algorithms
// It defines methods for adding values and computing a result based on those values.
type Computer interface {
	Push(float32)
	Value() float32
	PushAndValue(float32) float32
}

func computeChan(valueCh chan ChannelValue, computer Computer) chan ChannelValue {
	resCh := make(chan ChannelValue)

	go func() {
		for v := range valueCh {
			if v.ErrorMessage != "" {
				resCh <- v
				continue
			}

			q := computer.PushAndValue(v.Value)
			resCh <- ChannelValue{
				Value: q,
			}
		}

		close(resCh)
	}()

	return resCh
}
